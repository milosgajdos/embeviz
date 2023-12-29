package http

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/memory"
)

func MustProvidersService(t *testing.T, db *memory.DB) v1.ProvidersService {
	ps, err := memory.NewProvidersService(db)
	if err != nil {
		t.Fatal(err)
	}
	return ps
}

func MustServer(t *testing.T) *Server {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("failed to created new server: %v", err)
	}
	return s
}

func MustOpenDB(t *testing.T, dsn string) *memory.DB {
	db, err := memory.NewDB(dsn)
	if err != nil {
		t.Fatalf("failed creating new DB: %v", err)
	}
	if err := db.Open(); err != nil {
		t.Fatalf("failed opening DB: %v", err)
	}
	return db
}

func MustSeedProviders(t *testing.T, ps v1.ProvidersService, count int) []*v1.Provider {
	px := []*v1.Provider{}
	for i := 0; i < count; i++ {
		p, err := ps.AddProvider(context.TODO(), fmt.Sprintf("p%d", i), map[string]any{"foo": i})
		if err != nil {
			t.Fatal(err)
		}
		px = append(px, p)
	}
	return px
}

func MustSeedProviderEmbeddings(t *testing.T, db *memory.DB, p *v1.Provider, count int) {
	embs := []v1.Embedding{}
	for i := 0; i < count; i++ {
		embs = append(embs, v1.Embedding{})
	}
	store := map[string]map[string]any{
		p.UID: {
			"meta": p,
			"embs": embs,
			"proj": map[v1.Dim][]v1.Embedding{
				v1.Dim2D: embs,
				v1.Dim3D: embs,
			},
		},
	}
	if err := db.InitStore(store); err != nil {
		t.Fatalf("failed to init in-memory store: %v", err)
	}
}
