package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"client/client" 
)

func main() {
	fmt.Println("--- Uruchamiam Klienta API ---")

	baseURL := "https://jsonplaceholder.typicode.com"
	userAgent := "MójWłasnyKlientStudencki/1.0"

	api := client.NewAPIClient(baseURL, userAgent)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// --- 1. Testujemy GET ---
	fmt.Println("\nPobieranie posta o ID = 1...")
	post, err := api.GetPost(ctx, 1)
	if err != nil {
		log.Fatalf("Błąd podczas pobierania posta: %v", err)
	}
	fmt.Printf("Sukces! Pobrany Post:\n Tytuł: %s\n Treść: %s\n", post.Title, post.Body)

	// --- 2. Testujemy POST ---
	fmt.Println("\nTworzenie nowego posta...")
	newPost := &client.Post{
		Title:  "Nauka Golanga",
		Body:   "Tworzenie klienta HTTP krok po kroku.",
		UserID: 42,
	}

	createdPost, err := api.CreatePost(ctx, newPost)
	if err != nil {
		log.Fatalf("Błąd podczas tworzenia posta: %v", err)
	}
	fmt.Printf("Sukces! Serwer zwrócił utworzony obiekt (otrzymał ID: %d)\n", createdPost.ID)
}