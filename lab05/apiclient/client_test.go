package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sprawdzenie metody i ścieżki
		if r.Method != http.MethodGet {
			t.Errorf("oczekiwano GET, otrzymano %s", r.Method)
		}
		if r.URL.Path != "/items" {
			t.Errorf("oczekiwano /items, otrzymano %s", r.URL.Path)
		}
		
		// Sprawdzenie niestandardowego nagłówka User-Agent (Wymóg 5)
		if ua := r.Header.Get("User-Agent"); ua != "go-http-clientPV" {
			t.Errorf("oczekiwano User-Agent: go-http-clientPV, otrzymano %s", ua)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"First item","description":"Example description"}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	items, err := client.GetItems(context.Background())
	if err != nil {
		t.Fatalf("nieoczekiwany błąd: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("oczekiwano 1 elementu, otrzymano %d", len(items))
	}
	if items[0].Name != "First item" {
		t.Errorf("niepoprawna nazwa elementu: %s", items[0].Name)
	}
}

func TestCreateItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("oczekiwano POST, otrzymano %s", r.Method)
		}
		if r.URL.Path != "/items" {
			t.Errorf("oczekiwano /items, otrzymano %s", r.URL.Path)
		}
		
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("oczekiwano Content-Type application/json, otrzymano %s", ct)
		}

		var req CreateItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("błąd podczas dekodowania ciała żądania: %v", err)
		}
		if req.Name != "New item" {
			t.Errorf("oczekiwano nazwy 'New item', otrzymano %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":3,"name":"New item","description":"Description of the new item"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req := CreateItemRequest{
		Name:        "New item",
		Description: "Description of the new item",
	}

	item, err := client.CreateItem(context.Background(), req)
	if err != nil {
		t.Fatalf("nieoczekiwany błąd: %v", err)
	}

	if item == nil || item.ID != 3 {
		t.Errorf("niepoprawny zdekodowany obiekt: %+v", item)
	}
}

func TestUnexpectedStatusCode(t *testing.T) {
	// Symulacja błędu serwera 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	
	// Test dla GetItems
	_, err := client.GetItems(context.Background())
	if err == nil {
		t.Error("oczekiwano błędu dla niepoprawnego statusu, nie otrzymano żadnego")
	} else if err.Error() != "unexpected status code: 500" {
		t.Errorf("oczekiwano błędu o treści 'unexpected status code: 500', otrzymano %v", err)
	}
}