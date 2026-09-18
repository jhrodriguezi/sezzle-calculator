package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORS(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := WithCORS(next)

	t.Run("sets CORS headers and forwards non-preflight requests", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodPost, "/api/calculate", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected next handler to be called")
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "*")
		}
	})

	t.Run("short-circuits preflight OPTIONS requests", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodOptions, "/api/calculate", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if called {
			t.Fatal("expected next handler not to be called for OPTIONS")
		}
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})
}

func TestWithLogging(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := WithLogging(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
