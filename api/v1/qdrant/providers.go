package qdrant

import (
	"context"

	v1 "github.com/milosgajdos/embeviz/api/v1"
)

// ProvidersService allows to store data in qdrant vector store
type ProvidersService struct {
	db *DB
}

// NewProvidersService creates an instance of ProvidersService and returns it.
func NewProvidersService(db *DB) (*ProvidersService, error) {
	return &ProvidersService{
		db: db,
	}, nil
}

// AddProvider creates a new provider and returns it.
// It creates a new collection with the same name as the provider.
// nolint:revive
func (p *ProvidersService) AddProvider(ctx context.Context, name string, md map[string]any) (*v1.Provider, error) {
	return nil, v1.Errorf(v1.ENOTIMPLEMENTED, "AddProvider")
}

// GetProviders returns a list of providers filtered by filter.
// nolint:revive
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, int, error) {
	return nil, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviders")
}

// GetProviderByUID returns the provider with the given uuid.
// nolint:revive
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	return nil, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderByUID")
}

// GetProviderEmbeddings returns embeddings for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) GetProviderEmbeddings(ctx context.Context, uid string, filter v1.ProviderFilter) ([]v1.Embedding, int, error) {
	return nil, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderEmbeddings")
}

// GetProviderProjections returns embeddings projections for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) GetProviderProjections(ctx context.Context, uid string, filter v1.ProviderFilter) (map[v1.Dim][]v1.Embedding, int, error) {
	return nil, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderProjections")
}

// UpdateProviderEmbeddings generates embeddings for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, update v1.Embedding, proj v1.Projection) (*v1.Embedding, error) {
	return nil, v1.Errorf(v1.ENOTIMPLEMENTED, "UpdateProviderEmbeddings")
}

// DropProviderEmbeddings drops all provider embeddings from the store
// nolint:revive
func (p *ProvidersService) DropProviderEmbeddings(ctx context.Context, uid string) error {
	return v1.Errorf(v1.ENOTIMPLEMENTED, "DropProviderEmbeddings")
}

// ComputeProviderProjections drops existing projections and recomputes anew.
// nolint:revive
func (p *ProvidersService) ComputeProviderProjections(ctx context.Context, uid string, proj v1.Projection) error {
	return v1.Errorf(v1.ENOTIMPLEMENTED, "ComputeProviderProjections")
}
