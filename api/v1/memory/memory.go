package memory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sync"
)

const (
	// DSN is the in-memory data source name.
	DSN = ":memory:"
)

// Provider store model.
type Provider map[string]any

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

// Open opens the database connection.
func (db *DB) Open() (err error) {
	db.Lock()
	defer db.Unlock()

	if !db.Closed {
		return nil
	}

	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.DSN != DSN {
		if db.store, err = openFromFS(os.DirFS(db.DSN)); err != nil {
			return err
		}
		db.Closed = false
		return nil
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

// openFromFS opens DB and loads all data stored on the given fs.
func openFromFS(sys fs.FS) (map[string]map[string]any, error) {
	db := make(map[string]map[string]any)

	if err := fs.WalkDir(sys, ".", func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if e.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(sys, path)
		if err != nil {
			return err
		}

		p := make(Provider)
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
