package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Get_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("appId") != "test-app-id" {
			t.Errorf("appId = %q, want %q", r.URL.Query().Get("appId"), "test-app-id")
		}
		if r.URL.Query().Get("lang") != "J" {
			t.Errorf("lang = %q, want %q", r.URL.Query().Get("lang"), "J")
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-app-id", "J")
	var result map[string]string
	if err := client.Get("/json/getStatsList", nil, &result); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("status = %q, want %q", result["status"], "ok")
	}
}

func TestClient_Get_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "J")
	var result map[string]string
	err := client.Get("/json/getStatsList", nil, &result)
	if err == nil {
		t.Fatal("Get() should return error for 400 status")
	}
}

func TestClient_Get_networkError(t *testing.T) {
	client := NewClient("http://localhost:1", "test", "J")
	var result map[string]string
	err := client.Get("/json/getStatsList", nil, &result)
	if err == nil {
		t.Fatal("Get() should return error for bad host")
	}
}
