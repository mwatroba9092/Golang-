package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	// Importujemy nasz lokalny moduł
	"apiclient"
)

func main() {
	// 1. Uruchamiamy symulację zewnętrznego API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/items" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":1,"name":"First item","description":"Example description"}]`))
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/items" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id":3,"name":"New item","description":"Description of the new item"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// 2. Inicjalizujemy naszego klienta (z pakietu apiclient)
	client := apiclient.NewClient(server.URL)

	// 3. Tworzymy kontekst z limitem czasu
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 4. Testujemy metodę GET
	fmt.Println("--- Wykonuję żądanie GET ---")
	items, err := client.GetItems(ctx)
	if err != nil {
		log.Fatalf("Błąd pobierania elementów: %v", err)
	}
	fmt.Printf("Pobrane elementy: %+v\n\n", items)

	// 5. Testujemy metodę POST
	fmt.Println("--- Wykonuję żądanie POST ---")
	req := apiclient.CreateItemRequest{
		Name:        "New item",
		Description: "Description of the new item",
	}
	newItem, err := client.CreateItem(ctx, req)
	if err != nil {
		log.Fatalf("Błąd tworzenia elementu: %v", err)
	}
	fmt.Printf("Utworzony element: %+v\n", newItem)
}