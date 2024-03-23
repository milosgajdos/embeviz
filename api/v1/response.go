package v1

const (
	// DefaultLimit defines default results limit
	DefaultLimit = 20
)

// Page is used for API paging.
// Some upstream providers do not
// provide full count, but instead
// the ID of the next Page.
type Page struct {
	// Next is either a number
	// or a string ID which allows
	// resuming paging if provided.
	Next *string `json:"next,omitempty"`
	// Count is the number of all
	// results if provided.
	Count *int `json:"count,omitempty"`
}

// ChunkingResponse is returned when input chunking has been requested.
type ChunkingResponse struct {
	// Chunks contain indices into the chunked input.
	Chunks [][]int `json:"chunks"`
}

// ProvidersResponse is returned when querying providers.
type ProvidersResponse struct {
	Providers []*Provider `json:"providers"`
	Page      Page        `json:"page"`
}

// ProjectionsResponse is returned when querying provider embeddings projections
type ProjectionsResponse struct {
	Projections map[Dim][]Embedding `json:"embeddings"`
	Page        Page                `json:"page"`
}

// EmbeddingsResponse is returned when querying provider embeddings.
type EmbeddingsResponse struct {
	Embeddings []Embedding `json:"embeddings"`
	Page       Page        `json:"page"`
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}
