package memory

import "testing"

func MustProvidersService(t *testing.T, dsn string) *ProvidersService {
	db := MustOpenDB(t, dsn)
	ps, err := NewProvidersService(db)
	if err != nil {
		t.Fatal(err)
	}
	return ps
}

func MustClosedProvidersService(t *testing.T, dsn string) *ProvidersService {
	db := MustDB(t, dsn)
	ps, err := NewProvidersService(db)
	if err != nil {
		t.Fatal(err)
	}
	return ps
}

func MustDB(t *testing.T, dsn string) *DB {
	db, err := NewDB(dsn)
	if err != nil {
		t.Fatalf("failed creating new DB: %v", err)
	}
	return db
}

func MustOpenDB(t *testing.T, dsn string) *DB {
	db := MustDB(t, dsn)
	if err := db.Open(); err != nil {
		t.Fatalf("failed opening DB: %v", err)
	}
	return db
}
