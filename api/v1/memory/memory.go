package memory

import (
	"errors"
	"sync"
)

const (
	// DSN is the in-memory data source name.
	DSN = ":memory:"
)

var (
	ErrMissingDSN = errors.New("ErrMissingDSN")
	ErrInvalidDSN = errors.New("ErrInvalidDSN")
	ErrDBClosed   = errors.New("ErrDBClosed")
)

// DB is an in-memory DB store.
type DB struct {
	// DSN string.
	DSN string
	// Closed flag for DB operations.
	Closed bool
	// store stores memory store
	// TODO: maybe define some type
	// for the second level store
	store map[string]map[string]any
	*sync.RWMutex
}

// NewDB creates a new DB and returns it.
func NewDB(dsn string) (*DB, error) {
	return &DB{
		DSN:     dsn,
		Closed:  true,
		RWMutex: &sync.RWMutex{},
	}, nil
}

// InitStore lets you init the in-memory DB
// with the given store.
func (db *DB) InitStore(store map[string]map[string]any) error {
	db.Lock()
	defer db.Unlock()

	if err := db.open(); err != nil {
		return err
	}
	db.store = store

	return nil
}

// Open opens the database connection.
func (db *DB) Open() (err error) {
	db.Lock()
	defer db.Unlock()

	return db.open()
}

func (db *DB) open() error {
	if !db.Closed {
		return nil
	}

	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return ErrMissingDSN
	}

	if db.DSN != DSN {
		db.Closed = false
		return ErrInvalidDSN
	}

	db.store = make(map[string]map[string]any)
	db.Closed = false

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	db.Lock()
	defer db.Unlock()

	if db.Closed {
		return nil
	}

	db.Closed = true
	return nil
}
