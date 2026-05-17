package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchJSONSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"demo"}`))
	}))
	defer server.Close()

	var repo Repo
	statusCode, err := fetchJSON("token", server.URL, &repo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
	}
	if repo.Name != "demo" {
		t.Fatalf("expected repo name demo, got %q", repo.Name)
	}
}

func TestFetchJSONNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	var repo Repo
	statusCode, err := fetchJSON("token", server.URL, &repo)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
	}
}
