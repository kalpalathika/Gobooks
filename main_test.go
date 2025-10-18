package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- helpers ---

// resetBooks resets the global books slice before each test.
func resetBooks(seed []Book) {
	books = make([]Book, len(seed))
	copy(books, seed)
}

// doRequest sends an HTTP request through the router and returns the recorder.
func doRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &buf)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// apiResponse is a typed version of your ApiResponse for decoding Data.
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Error   string `json:"error"`
}

// --- tests ---

func TestGetBooks_OK(t *testing.T) {
	// Arrange
	seed := []Book{
		{ID: "1", Title: "Book1", Author: "Author1", Year: 1997},
		{ID: "2", Title: "Book2", Author: "Author2", Year: 1998},
	}
	resetBooks(seed)
	router := newRouter()

	// Act
	rec := doRequest(t, router, http.MethodGet, "/books", nil)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", rec.Code, rec.Body.String())
	}

	var resp apiResponse[[]Book]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v\nbody=%s", err, rec.Body.String())
	}

	if !resp.Success {
		t.Fatalf("expected success=true, got false; error=%s", resp.Error)
	}
	if len(resp.Data) != len(seed) {
		t.Fatalf("expected %d books, got %d", len(seed), len(resp.Data))
	}
	if resp.Data[0].Title != seed[0].Title {
		t.Fatalf("expected first title %q, got %q", seed[0].Title, resp.Data[0].Title)
	}
}
