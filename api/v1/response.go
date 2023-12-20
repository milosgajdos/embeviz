package v1

const (
	// DefaultLimit defines default results limit
	DefaultLimit = 20
)

// ProvidersResponse is returned when querying providers.
type ProvidersResponse struct {
	Providers []*Provider `json:"providers"`
	N         int         `json:"n"`
}

// ProjectionsResponse is returned when querying provider embeddings projections
type ProjectionsResponse struct {
	Embeddings map[Dim][]Embedding `json:"embeddings"`
	N          int                 `json:"n"`
}

// EmbeddingsResponse is returned when querying provider embeddings.
type EmbeddingsResponse struct {
	Embeddings []Embedding `json:"embeddings"`
	N          int         `json:"n"`
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}
