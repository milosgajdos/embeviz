package qdrant

import (
	"crypto/tls"
	"errors"
	"os"
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
	// metadata
	md metadata.MD
	// TODO: here to handle API endpoints
	// which are not yet supported by gRPC
	httpClient *HTTPClient
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
	if db.DSN == "" {
		return ErrMissingDSN
	}

	scheme, authKey, hostAddr, err := parseDSN(db.DSN)
	if err != nil {
		return err
	}

	// TODO: build HTTP API Base URL here and pass it to the HTTP API client.
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

	// NOTE: we expect this to be an API key for the qdrant cloud.
	// if it's set we update the connection context with auth details.
	if authKey == "" {
		authKey = os.Getenv("QDRANT_API_KEY")
	}
	if authKey != "" {
		db.md = metadata.New(map[string]string{
			"api-key": authKey,
		})
	}

	conn, err := grpc.Dial(hostAddr, dialOpts...)
	if err != nil {
		return err
	}
	db.conn = conn
	db.col = pb.NewCollectionsClient(conn)
	db.pts = pb.NewPointsClient(conn)

	db.httpClient = NewHTTPClient() // TODO: pass options

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Close database.
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func parseDSN(dsn string) (scheme, authKey, hostAddr string, err error) {
	var (
		ok   bool
		path string
	)

	scheme, path, ok = strings.Cut(dsn, "://")
	if !ok {
		err = ErrInvalidDSN
		return
	}

	authKey, hostAddr, ok = strings.Cut(path, "@")
	if !ok {
		err = ErrInvalidDSN
		return
	}

	// NOTE: maybe we could init this to localhost:6334 by default
	// instead of returning error
	if hostAddr == "" {
		err = ErrInvalidDSN
		return
	}
	return
}
