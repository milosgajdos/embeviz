package memory

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	v1 "github.com/milosgajdos/embeviz/api/v1"
)

func TestCreateGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		name := "foo"
		md := map[string]any{
			"foo": "bar",
		}

		p, err := ps.AddProvider(context.TODO(), name, md)
		if err != nil {
			t.Fatal(err)
		}

		if p.UID == "" {
			t.Fatal("expected non-emnpty UID")
		}
		if p.Name != name {
			t.Fatalf("expected name: %v, got: %v", name, p.Name)
		}
		if !reflect.DeepEqual(p.Metadata, md) {
			t.Error("invalid provider metadata")
		}
	})

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		name := "foo"
		md := map[string]any{
			"foo": "bar",
		}

		if _, err := ps.AddProvider(context.TODO(), name, md); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestGetProviders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ps := MustProvidersService(t, DSN)

	// insert test data
	pCount := 5
	for i := 0; i < pCount; i++ {
		if _, err := ps.AddProvider(context.TODO(), fmt.Sprintf("p%d", i), map[string]any{"foo": i}); err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		name       string
		filter     v1.ProviderFilter
		expRes     int
		expMatches int
		expErr     bool
	}{
		{"All", v1.ProviderFilter{}, pCount, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 100}, 0, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: -1}, pCount, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1}, pCount - 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 1}, 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 10}, pCount - 1, pCount, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			px, page, err := ps.GetProviders(context.TODO(), tc.filter)
			if !tc.expErr && err != nil {
				t.Fatal(err)
			}

			if *page.Count != tc.expMatches {
				t.Errorf("expected providers: %d, got: %d", tc.expMatches, *page.Count)
			}

			if tc.expRes != len(px) {
				t.Errorf("expected results: %d, got: %d", tc.expRes, len(px))
			}
		})
	}

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, _, err := ps.GetProviders(context.TODO(), v1.ProviderFilter{}); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestGetProviderByUID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		name := "foo"
		md := map[string]any{
			"foo": "bar",
		}

		p, err := ps.AddProvider(context.TODO(), name, md)
		if err != nil {
			t.Fatal(err)
		}

		px, err := ps.GetProviderByUID(context.TODO(), p.UID)
		if err != nil {
			t.Fatal(err)
		}

		if p.UID != px.UID {
			t.Fatalf("Expected UID: %s, got: %s", p.UID, px.UID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		if _, err := ps.GetProviderByUID(context.TODO(), "garbageUID"); v1.ErrorCode(err) != v1.ENOTFOUND {
			t.Fatalf("expected error: %s, got: %s", v1.ENOTFOUND, v1.ErrorCode(err))
		}
	})

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, err := ps.GetProviderByUID(context.TODO(), "garbageUID"); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestGetProviderEmbeddings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ps := MustProvidersService(t, DSN)

	name := "foo"
	md := map[string]any{
		"foo": "bar",
	}

	p, err := ps.AddProvider(context.TODO(), name, md)
	if err != nil {
		t.Fatal(err)
	}

	// insert test data
	embs := []v1.Embedding{}
	pCount := 5
	for i := 0; i < pCount; i++ {
		embs = append(embs, v1.Embedding{})
	}
	ps.db.store[p.UID][emb] = embs

	testCases := []struct {
		name       string
		filter     v1.ProviderFilter
		expRes     int
		expMatches int
		expErr     bool
	}{
		{"All", v1.ProviderFilter{}, pCount, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 100}, 0, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: -1}, pCount, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1}, pCount - 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 1}, 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 10}, pCount - 1, pCount, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			px, page, err := ps.GetProviderEmbeddings(context.TODO(), p.UID, tc.filter)
			if !tc.expErr && err != nil {
				t.Fatal(err)
			}

			if *page.Count != tc.expMatches {
				t.Errorf("expected embeddings: %d, got: %d", tc.expMatches, *page.Count)
			}

			if tc.expRes != len(px) {
				t.Errorf("expected results: %d, got: %d", tc.expRes, len(px))
			}
		})
	}

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, _, err := ps.GetProviders(context.TODO(), v1.ProviderFilter{}); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestGetProviderProjections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ps := MustProvidersService(t, DSN)

	name := "foo"
	md := map[string]any{
		"foo": "bar",
	}

	p, err := ps.AddProvider(context.TODO(), name, md)
	if err != nil {
		t.Fatal(err)
	}

	// insert test data
	// embeddings
	embs := []v1.Embedding{}
	pCount := 5
	for i := 0; i < pCount; i++ {
		embs = append(embs, v1.Embedding{})
	}
	ps.db.store[p.UID][emb] = embs
	// projections
	proj2D := []v1.Embedding{}
	proj3D := []v1.Embedding{}
	for i := 0; i < pCount; i++ {
		proj2D = append(proj2D, v1.Embedding{})
		proj3D = append(proj3D, v1.Embedding{})
	}
	ps.db.store[p.UID][proj] = map[v1.Dim][]v1.Embedding{
		v1.Dim2D: proj2D,
		v1.Dim3D: proj3D,
	}

	dim2D := v1.Dim2D
	//dim3D := v1.Dim3D

	testCases := []struct {
		name       string
		filter     v1.ProviderFilter
		expRes     int
		expMatches int
		expErr     bool
	}{
		{"All", v1.ProviderFilter{}, pCount, pCount, false},
		{"AllWithDim", v1.ProviderFilter{Dim: &dim2D}, pCount, pCount, false},
		{"LimitWithDim", v1.ProviderFilter{Limit: 1, Dim: &dim2D}, 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 100}, 0, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: -1}, pCount, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1}, pCount - 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 1}, 1, pCount, false},
		{"LimitOffset", v1.ProviderFilter{Offset: 1, Limit: 10}, pCount - 1, pCount, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			px, page, err := ps.GetProviderProjections(context.TODO(), p.UID, tc.filter)
			if !tc.expErr && err != nil {
				t.Fatal(err)
			}

			if *page.Count != tc.expMatches {
				t.Errorf("expected projections: %d, got: %d", tc.expMatches, *page.Count)
			}

			// NOTE: we seed the data so we expect the filter seeds to be retrurned
			// and they are guaranteed to have some projections
			pxRes := len(px[v1.Dim2D])

			if tc.expRes != pxRes {
				t.Errorf("expected results: %d, got: %d", tc.expRes, pxRes)
			}
		})
	}

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, _, err := ps.GetProviders(context.TODO(), v1.ProviderFilter{}); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestUpdateProviderEmbeddings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		name := "foo"
		md := map[string]any{
			"foo": "bar",
		}

		p, err := ps.AddProvider(context.TODO(), name, md)
		if err != nil {
			t.Fatal(err)
		}

		emb := v1.Embedding{
			Values: []float64{1.0, 2.0, 3.0, 4.0},
		}

		e, err := ps.UpdateProviderEmbeddings(context.TODO(), p.UID, emb, v1.PCA)
		if err != nil {
			t.Fatalf("expected error: %s", err)
		}

		if !reflect.DeepEqual(e.Values, emb.Values) {
			t.Fatalf("expected: %v, got: %v", emb.Values, e.Values)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		emb := v1.Embedding{}

		if _, err := ps.UpdateProviderEmbeddings(context.TODO(), "garbageUID", emb, v1.PCA); v1.ErrorCode(err) != v1.ENOTFOUND {
			t.Fatalf("expected error: %s, got: %s", v1.ENOTFOUND, v1.ErrorCode(err))
		}
	})

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, err := ps.GetProviderByUID(context.TODO(), "garbageUID"); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}

func TestDropProviderEmbeddings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		name := "foo"
		md := map[string]any{
			"foo": "bar",
		}

		p, err := ps.AddProvider(context.TODO(), name, md)
		if err != nil {
			t.Fatal(err)
		}

		if err := ps.DropProviderEmbeddings(context.TODO(), p.UID); err != nil {
			t.Fatalf("expected error: %s", err)
		}

		px, page, err := ps.GetProviderProjections(context.TODO(), p.UID, v1.ProviderFilter{})
		if err != nil {
			t.Fatal(err)
		}

		if dim2D := len(px[v1.Dim2D]); dim2D != 0 {
			t.Fatalf("expected no projections, got: %d", dim2D)
		}

		if dim3D := len(px[v1.Dim3D]); dim3D != 0 {
			t.Fatalf("expected no projections, got: %d", dim3D)
		}

		if *page.Count != 0 {
			t.Fatalf("expected no projections, got: %d", page)
		}

		ex, page, err := ps.GetProviderEmbeddings(context.TODO(), p.UID, v1.ProviderFilter{})
		if err != nil {
			t.Fatal(err)
		}

		if len(ex) != 0 {
			t.Fatalf("expected no projections, got: %d", len(ex))
		}

		if *page.Count != 0 {
			t.Fatalf("expected no projections, got: %d", page)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ps := MustProvidersService(t, DSN)

		if err := ps.DropProviderEmbeddings(context.TODO(), "blahUID"); v1.ErrorCode(err) != v1.ENOTFOUND {
			t.Fatalf("expected error: %s, got: %s", v1.ENOTFOUND, v1.ErrorCode(err))
		}
	})

	t.Run("ClosedDB", func(t *testing.T) {
		ps := MustClosedProvidersService(t, DSN)

		if _, err := ps.GetProviderByUID(context.TODO(), "garbageUID"); !errors.Is(err, ErrDBClosed) {
			t.Fatalf("expected error: %v, got: %v", ErrDBClosed, err)
		}
	})
}
