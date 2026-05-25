package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPost_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedUserAgent := "MojaAplikacja/1.0"
		if r.Header.Get("User-Agent") != expectedUserAgent {
			t.Errorf("Oczekiwano User-Agent %q, otrzymano %q", expectedUserAgent, r.Header.Get("User-Agent"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 42, "title": "Testowe API", "body": "Działa!", "userId": 1}`))
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "MojaAplikacja/1.0")

	post, err := client.GetPost(context.Background(), 42)
	
	if err != nil {
		t.Fatalf("Oczekiwano braku błędów, otrzymano: %v", err)
	}

	if post.ID != 42 || post.Title != "Testowe API" {
		t.Errorf("Zdekodowano niepoprawne dane: %+v", post)
	}
}

func TestGetPost_NotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "MojaAplikacja/1.0")

	_, err := client.GetPost(context.Background(), 999)
	if err == nil {
		t.Fatal("Oczekiwano błędu z powodu statusu 404, jednak błąd nie wystąpił")
	}
}