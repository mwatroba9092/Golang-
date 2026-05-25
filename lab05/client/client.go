package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Post struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type userAgentTransport struct {
	transport http.RoundTripper
	userAgent string
}

func (c *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := c.transport
	if rt == nil {
		rt = http.DefaultTransport
	}

	reqClone := req.Clone(req.Context())
	reqClone.Header.Set("User-Agent", c.userAgent)

	return rt.RoundTrip(reqClone)
}

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL, userAgent string) *APIClient {
	customTransport := &userAgentTransport{
		transport: http.DefaultTransport,
		userAgent: userAgent,
	}

	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: customTransport,
			Timeout:   15 * time.Second,
		},
	}
}

func (c *APIClient) GetPost(ctx context.Context, id int) (*Post, error) {
	url := fmt.Sprintf("%s/posts/%d", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("błąd tworzenia żądania: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd sieci podczas żądania: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("post o ID %d nie istnieje (404)", id)
		}
		return nil, fmt.Errorf("błąd API, oczekiwano 200 OK, otrzymano: %d", resp.StatusCode)
	}

	var post Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return &post, nil
}

func (c *APIClient) CreatePost(ctx context.Context, newPost *Post) (*Post, error) {
	url := fmt.Sprintf("%s/posts", c.baseURL)

	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(newPost); err != nil {
		return nil, fmt.Errorf("błąd kodowania danych: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("błąd tworzenia żądania: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nieoczekiwany kod odpowiedzi: %d", resp.StatusCode)
	}

	var createdPost Post
	if err := json.NewDecoder(resp.Body).Decode(&createdPost); err != nil {
		return nil, fmt.Errorf("błąd dekodowania zwrotnego JSON: %w", err)
	}

	return &createdPost, nil
}