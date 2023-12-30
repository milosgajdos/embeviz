package qdrant

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	ErrMissingDSN = errors.New("ErrMissingDSN")
	ErrInvalidDSN = errors.New("ErrInvalidDSN")
	ErrDBClosed   = errors.New("ErrDBClosed")
)

// DB is qdrant DB store handle
type DB struct {
	// DSN string
	DSN string
	// conn maintains connetion
	conn *grpc.ClientConn
	// collections client
	col pb.CollectionsClient
	// points client
	pts pb.PointsClient
	// auth context
	ctx    context.Context // background context
	cancel func()          // cancel background context
}

// NewDB creates a new DB and returns it.
// TODO: figure out the DSN parsing,
// For now we support a very silly DSN parsing
// and expect the DSN to look like thi:
// qadrnt(s)://secret@host
func NewDB(dsn string) (*DB, error) {
	return &DB{
		DSN: dsn,
	}, nil
}

// Open opens the database connection.
func (db *DB) Open() (err error) {
	scheme, path, ok := strings.Cut(db.DSN, "://")
	if !ok {
		return ErrInvalidDSN
	}
	var dialOpts []grpc.DialOption
	switch scheme {
	case "qdrant://":
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	case "qdrants://":
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	default:
		// NOTE: we default to insecure opts
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	authKey, hostAddr, ok := strings.Cut(path, "@")
	if !ok {
		return ErrInvalidDSN
	}

	// NOTE: we could support the default localhost here: localhost:6334
	if hostAddr == "" {
		return ErrInvalidDSN
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())

	// NOTE: we expect this to be an API key for qdrant cloud.
	// if it's set we update the connection context with auth details.
	if authKey != "" {
		md := metadata.New(map[string]string{"api-key": authKey})
		db.ctx = metadata.NewOutgoingContext(db.ctx, md)
	}

	conn, err := grpc.Dial(hostAddr, dialOpts...)
	if err != nil {
		return err
	}
	db.conn = conn
	db.col = pb.NewCollectionsClient(conn)
	db.pts = pb.NewPointsClient(conn)

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
