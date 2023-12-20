package memory

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/internal/projection"
)

const (
	meta = "meta"
	emb  = "embs"
	proj = "proj"
)

type ProvidersService struct {
	db *DB
}

// NewProvidersService creates an instance of ProvidersService and returns it.
func NewProvidersService(db *DB) (*ProvidersService, error) {
	return &ProvidersService{
		db: db,
	}, nil
}

// AddProvider adds a new embeddings provider.
// nolint:revive
func (p *ProvidersService) AddProvider(ctx context.Context, name string, md map[string]any) (*v1.Provider, error) {
	p.db.Lock()
	defer p.db.Unlock()
	uid := uuid.New().String()
	provider := &v1.Provider{
		UID:      uid,
		Name:     name,
		Metadata: md,
	}
	p.db.store[uid] = map[string]any{
		meta: provider,
		emb:  []v1.Embedding{},
		proj: map[v1.Dim][]v1.Embedding{
			v1.Dim2D: {},
			v1.Dim3D: {},
		},
	}
	return provider, nil
}

// GetProviders fetches all available providers.
// nolint:revive
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, int, error) {
	// TODO: pagination
	p.db.RLock()
	defer p.db.RUnlock()
	px := make([]*v1.Provider, 0, len(p.db.store[meta]))
	for _, p := range p.db.store {
		px = append(px, p[meta].(*v1.Provider))
	}
	return px, len(px), nil
}

// GetProviderByid fetches a specific provider by uid.
// nolint:revive
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	p.db.RLock()
	defer p.db.RUnlock()
	if p, ok := p.db.store[uid]; ok {
		return p[meta].(*v1.Provider), nil

	}
	return nil, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
}

// GetProviderEmbeddings fetches a specific provider embeddings.
// nolint:revive
func (p *ProvidersService) GetProviderEmbeddings(ctx context.Context, uid string, filter v1.ProviderFilter) ([]v1.Embedding, int, error) {
	// TODO: pagination
	p.db.RLock()
	defer p.db.RUnlock()
	provider, ok := p.db.store[uid]
	if !ok {
		return nil, 0, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
	}
	embs := provider[emb].([]v1.Embedding)
	newEmbs := make([]v1.Embedding, len(embs))
	copy(newEmbs, embs)
	return newEmbs, len(newEmbs), nil
}

// GetProviderProjections fetches a specific provider projection.
// nolint:revive
func (p *ProvidersService) GetProviderProjections(ctx context.Context, uid string, filter v1.ProviderFilter) (map[v1.Dim][]v1.Embedding, int, error) {
	// TODO: pagination
	p.db.RLock()
	defer p.db.RUnlock()
	provider, ok := p.db.store[uid]
	if !ok {
		return nil, 0, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
	}
	if dim := filter.Dim; dim != nil {
		if *dim != v1.Dim2D && *dim != v1.Dim3D {
			return nil, 0, v1.Errorf(v1.EINVALID, "invalid dimension %v for provider %q", *dim, uid)
		}
		projStore := provider[proj].(map[v1.Dim][]v1.Embedding)
		projections := projStore[*dim]
		newProjections := make([]v1.Embedding, len(projections))
		copy(newProjections, projections)
		return map[v1.Dim][]v1.Embedding{*filter.Dim: newProjections}, len(projections), nil
	}
	// no filter returns all dimensions projectsion
	projections := provider[proj].(map[v1.Dim][]v1.Embedding)
	newProjections2D := make([]v1.Embedding, len(projections[v1.Dim2D]))
	copy(newProjections2D, projections[v1.Dim2D])
	newProjections3D := make([]v1.Embedding, len(projections[v1.Dim3D]))
	copy(newProjections3D, projections[v1.Dim3D])
	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: newProjections2D,
		v1.Dim3D: newProjections3D,
	}, len(newProjections2D), nil
}

// UpdateProviderEmbeddings updates embeddings of a specific provider.
// nolint:revive
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, embed v1.Embedding, prj v1.Projection) (*v1.Embedding, error) {
	p.db.Lock()
	defer p.db.Unlock()
	provider, ok := p.db.store[uid]
	if !ok {
		return nil, v1.Errorf(v1.ENOTFOUND, "provider %s not found", uid)
	}
	embs := provider[emb].([]v1.Embedding)
	newEmbs := make([]v1.Embedding, len(embs))
	copy(newEmbs, embs)
	newEmbs = append(newEmbs, embed)

	var (
		err    error
		proj2D []v1.Embedding
		proj3D []v1.Embedding
	)

	// Calculate projection
	switch prj {
	case v1.PCA:
		proj2D, err = projection.PCA(newEmbs, v1.Dim2D)
		if err != nil {
			return nil, err
		}
		proj3D, err = projection.PCA(newEmbs, v1.Dim3D)
		if err != nil {
			return nil, err
		}
	case v1.TSNE:
		proj2D, err = projection.TSNE(newEmbs, v1.Dim2D)
		if err != nil {
			return nil, err
		}
		proj3D, err = projection.TSNE(newEmbs, v1.Dim3D)
		if err != nil {
			return nil, err
		}
	default:
		return nil, v1.Errorf(v1.EINVALID, "invalid projection: %v", proj)

	}
	provider[emb] = newEmbs
	provider[proj] = map[v1.Dim][]v1.Embedding{
		v1.Dim2D: proj2D,
		v1.Dim3D: proj3D,
	}
	return &embed, nil
}
