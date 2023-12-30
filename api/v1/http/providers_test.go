package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/memory"
)

func TestGetAllProviders(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		count := 2
		_ = MustSeedProviders(t, ps, count)

		req := httptest.NewRequest("GET", "/api/v1/providers", nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		ret := new(v1.ProvidersResponse)
		if err := json.Unmarshal(body, ret); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *ret.Page.Count != count {
			t.Errorf("expected providers: %d, got: %d", count, *ret.Page.Count)
		}
	})

	t.Run("200/Paginated", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		limit, offset := 2, 2

		count := 5
		_ = MustSeedProviders(t, ps, count)

		urlPath := fmt.Sprintf("/api/v1/providers?limit=%d&offset=%d", limit, offset)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		ret := new(v1.ProvidersResponse)
		if err := json.Unmarshal(body, ret); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *ret.Page.Count != count {
			t.Errorf("expected total providers: %d, got: %d", count, *ret.Page.Count)
		}

		if n := len(ret.Providers); n != limit {
			t.Errorf("expected providers: %d, got: %d", limit, n)
		}
	})

	t.Run("500", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		gs := MustProvidersService(t, db)
		s.ProvidersService = gs

		// NOTE: we simulate the loss of DB connection like this.
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB: %v", err)
		}

		req := httptest.NewRequest("GET", "/api/v1/providers", nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected status code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}

func TestGetProviderByUID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		px := MustSeedProviders(t, ps, 2)

		uid := px[0].UID
		urlPath := fmt.Sprintf("/api/v1/providers/%s", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		p := new(v1.Provider)
		if err := json.Unmarshal(body, p); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if p.UID != uid {
			t.Fatalf("expected graph uid: %s, got: %s", uid, p.UID)
		}
	})

	t.Run("400", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "dflksdjfdlksf"
		urlPath := fmt.Sprintf("/api/v1/providers/%s", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusBadRequest {
			t.Fatalf("expected status code: %d, got: %d", http.StatusBadRequest, code)
		}
	})

	t.Run("404", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusNotFound {
			t.Fatalf("expected status code: %d, got: %d", http.StatusNotFound, code)
		}
	})

	t.Run("500", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		// we simulate the loss of DB connection like this.
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB: %v", err)
		}

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected status code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}

func TestGetProviderEmbeddings(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		px := MustSeedProviders(t, ps, 1)
		uid := px[0].UID

		count := 5
		MustSeedProviderEmbeddings(t, db, px[0], count)

		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		ret := new(v1.EmbeddingsResponse)
		if err := json.Unmarshal(body, ret); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *ret.Page.Count != count {
			t.Errorf("expected providers: %d, got: %d", count, *ret.Page.Count)
		}
	})

	t.Run("400", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "dflksdjfdlksf"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusBadRequest {
			t.Fatalf("expected status code: %d, got: %d", http.StatusBadRequest, code)
		}
	})

	t.Run("404", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusNotFound {
			t.Fatalf("expected status code: %d, got: %d", http.StatusNotFound, code)
		}
	})

	t.Run("500", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		// we simulate the loss of DB connection like this.
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB: %v", err)
		}

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected status code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}

func TestGetProviderProjections(t *testing.T) {
	t.Run("200/ALL", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		px := MustSeedProviders(t, ps, 1)
		uid := px[0].UID

		count := 5
		MustSeedProviderEmbeddings(t, db, px[0], count)

		urlPath := fmt.Sprintf("/api/v1/providers/%s/projections", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		ret := new(v1.ProjectionsResponse)
		if err := json.Unmarshal(body, ret); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *ret.Page.Count != count {
			t.Errorf("expected providers: %d, got: %d", count, *ret.Page.Count)
		}

		// we should get both 2D and 2D projections back
		if len(ret.Projections) != 2 {
			t.Errorf("expected both 2D and 3D projections, got: %d", len(ret.Projections))
		}
	})

	t.Run("200/2D_Only", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		px := MustSeedProviders(t, ps, 1)
		uid := px[0].UID

		count := 5
		MustSeedProviderEmbeddings(t, db, px[0], count)

		dim := "2d"

		urlPath := fmt.Sprintf("/api/v1/providers/%s/projections?dim=%s", uid, dim)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		ret := new(v1.ProjectionsResponse)
		if err := json.Unmarshal(body, ret); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *ret.Page.Count != count {
			t.Errorf("expected providers: %d, got: %d", count, *ret.Page.Count)
		}

		// we should only get 2D projections
		if len(ret.Projections) != 1 {
			t.Errorf("expected only %s projectionxs, got: %d", dim, len(ret.Projections))
		}
	})

	t.Run("400", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "dflksdjfdlksf"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/projections", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusBadRequest {
			t.Fatalf("expected status code: %d, got: %d", http.StatusBadRequest, code)
		}
	})

	t.Run("404", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/projections", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusNotFound {
			t.Fatalf("expected status code: %d, got: %d", http.StatusNotFound, code)
		}
	})

	t.Run("500", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		// we simulate the loss of DB connection like this.
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB: %v", err)
		}

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/projections", uid)
		req := httptest.NewRequest("GET", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected status code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}

func TestDropProviderEmbeddings(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		px := MustSeedProviders(t, ps, 1)
		uid := px[0].UID

		count := 5
		MustSeedProviderEmbeddings(t, db, px[0], count)
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("DELETE", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusNoContent {
			t.Fatalf("expected status code: %d, got: %d", http.StatusNoContent, code)
		}

		// Verify no embeddings are returned for this provider
		urlPath = fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req = httptest.NewRequest("GET", urlPath, nil)

		resp, err = s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		embRet := new(v1.EmbeddingsResponse)
		if err := json.Unmarshal(body, embRet); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *embRet.Page.Count != 0 {
			t.Errorf("expected providers: %d, got: %d", 0, *embRet.Page.Count)
		}

		// Verify no projections are returned for this provider
		urlPath = fmt.Sprintf("/api/v1/providers/%s/projections", uid)
		req = httptest.NewRequest("GET", urlPath, nil)

		resp, err = s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected status code: %d, got: %d", http.StatusOK, code)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		projRet := new(v1.ProjectionsResponse)
		if err := json.Unmarshal(body, projRet); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if *projRet.Page.Count != 0 {
			t.Errorf("expected providers: %d, got: %d", 0, *projRet.Page.Count)
		}
	})

	t.Run("400", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		uid := "sdlfkjsdflkdjf"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("DELETE", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusBadRequest {
			t.Fatalf("expected status code: %d, got: %d", http.StatusBadRequest, code)
		}
	})

	t.Run("500", func(t *testing.T) {
		s := MustServer(t)
		db := MustOpenDB(t, memory.DSN)
		ps := MustProvidersService(t, db)
		s.ProvidersService = ps

		// we simulate the loss of DB connection like this.
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB: %v", err)
		}

		uid := "97153afd-c434-4ca0-a35b-7467fcd08df1"
		urlPath := fmt.Sprintf("/api/v1/providers/%s/embeddings", uid)
		req := httptest.NewRequest("DELETE", urlPath, nil)

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected status code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}
