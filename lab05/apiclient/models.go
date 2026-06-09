package apiclient

// Item reprezentuje pojedynczy element zwracany przez API.
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateItemRequest reprezentuje dane wejściowe do utworzenia nowego elementu.
type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}