package v1

import "context"

// Provider for embeddings.
type Provider struct {
	// UID of the provider's UUID.
	UID string `json:"id"`
	// Name is the name of the provider
	Name string `json:"name"`
	// Metadata about the provider.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Embedding is vector embedding.
type Embedding struct {
	// Values stores embedding vector values.
	// NOTE: the key is set to value - singular
	// because the API is consumed by ECharts and
	// it's just sad ECharts expects value slice.
	// We could handle that in JS but who can be bothered?
	Values []float64 `json:"value"`
	// Metadata for the given embedding vector.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Dim is projection dimenstion
type Dim string

const (
	Dim2D Dim = "2D"
	Dim3D Dim = "3D"
)

// Projection algorithm.
type Projection string

const (
	// TSNE projection
	// https://en.wikipedia.org/wiki/T-distributed_stochastic_neighbor_embedding
	TSNE Projection = "tsne"
	// PCA projection
	// https://en.wikipedia.org/wiki/Principal_component_analysis
	PCA Projection = "pca"
)

// ProviderFilter is used for filtering providers.
type ProviderFilter struct {
	// Filtering fields.
	Dim *Dim `json:"dim"`
	// Restrict to subset of range.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// EmbeddingUpdate is used to fetch embeddings.
type EmbeddingUpdate struct {
	Text       string         `json:"text"`
	Label      string         `json:"label"`
	Projection Projection     `json:"projection"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// EmbeddingProjectionUpdate is used to recompute embedding projections.
type EmbeddingProjectionUpdate struct {
	Projection Projection     `json:"projection"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// ProvidersService manages embedding providers.
type ProvidersService interface {
	// AddProvider creates a new provider and returns it.
	AddProvider(ctx context.Context, name string, metadata map[string]any) (*Provider, error)
	// GetProviders returns a list of providers filtered by filter.
	GetProviders(ctx context.Context, filter ProviderFilter) ([]*Provider, int, error)
	// GetProviderByUID returns the provider with the given uuid.
	GetProviderByUID(ctx context.Context, uid string) (*Provider, error)
	// GetProviderEmbeddings returns embeddings for the provider with the given uid.
	GetProviderEmbeddings(ctx context.Context, uid string, filter ProviderFilter) ([]Embedding, int, error)
	// GetProviderProjections returns embeddings projections for the provider with the given uid.
	GetProviderProjections(ctx context.Context, uid string, filter ProviderFilter) (map[Dim][]Embedding, int, error)
	// UpdateProviderEmbeddings generates embeddings for the provider with the given uid.
	UpdateProviderEmbeddings(ctx context.Context, uid string, update Embedding, projection Projection) (*Embedding, error)
	// DropProviderEmbeddings drops all provider embeddings from the store
	DropProviderEmbeddings(ctx context.Context, uid string) error
	// ComputeProviderProjections drops existing projections and recomputes anew.
	ComputeProviderProjections(ctx context.Context, uid string, projection Projection) error
}
