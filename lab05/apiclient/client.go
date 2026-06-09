package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client to klient HTTP API z wielokrotnie używanym http.Client.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient tworzy nową instancję klienta API z wstrzykniętym nagłówkiem User-Agent.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: &UserAgentTransport{
				UserAgent: "go-http-clientPV",
			},
		},
	}
}

// GetItems pobiera listę elementów używając endpointu GET /items.
func (c *Client) GetItems(ctx context.Context) ([]Item, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/items", nil)
	if err != nil {
		return nil, fmt.Errorf("błąd tworzenia żądania: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("żądanie nie powiodło się: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return items, nil
}

// CreateItem tworzy nowy element używając endpointu POST /items.
func (c *Client) CreateItem(ctx context.Context, input CreateItemRequest) (*Item, error) {
	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("błąd kodowania wejścia do JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/items", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("błąd tworzenia żądania: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("żądanie nie powiodło się: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return &item, nil
}