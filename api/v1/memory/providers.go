package memory

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/internal/paging"
	"github.com/milosgajdos/embeviz/api/v1/internal/projection"
)

const (
	// metadata keyspace
	meta = "meta"
	// embeddings keyspace
	emb = "embs"
	// projection keyspace
	proj = "proj"
)

// ProvidersService is an in-memory store for embeddings providers.
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
	if p.db.Closed {
		return nil, v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

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
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, v1.Page, error) {
	p.db.RLock()
	defer p.db.RUnlock()
	count := 0
	if p.db.Closed {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	px := make([]*v1.Provider, 0, len(p.db.store[meta]))
	for _, p := range p.db.store {
		px = append(px, p[meta].(*v1.Provider))
	}
	count = len(px)
	offset, ok := filter.Offset.(int)
	if !ok {
		offset = 0
	}
	return paging.ApplyOffsetLimit(px, offset, filter.Limit).([]*v1.Provider), v1.Page{Count: &count}, nil
}

// GetProviderByid fetches a specific provider by uid.
// nolint:revive
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	p.db.RLock()
	defer p.db.RUnlock()
	if p.db.Closed {
		return nil, v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	if p, ok := p.db.store[uid]; ok {
		return p[meta].(*v1.Provider), nil

	}
	return nil, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
}

// GetProviderEmbeddings fetches a specific provider embeddings.
// nolint:revive
func (p *ProvidersService) GetProviderEmbeddings(ctx context.Context, uid string, filter v1.ProviderFilter) ([]v1.Embedding, v1.Page, error) {
	p.db.RLock()
	defer p.db.RUnlock()
	count := 0
	if p.db.Closed {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	provider, ok := p.db.store[uid]
	if !ok {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
	}
	embs := provider[emb].([]v1.Embedding)
	newEmbs := make([]v1.Embedding, len(embs))
	copy(newEmbs, embs)
	count = len(newEmbs)
	offset, ok := filter.Offset.(int)
	if !ok {
		offset = 0
	}
	return paging.ApplyOffsetLimit(newEmbs, offset, filter.Limit).([]v1.Embedding), v1.Page{Count: &count}, nil
}

// GetProviderProjections fetches a specific provider embeddings projection.
// nolint:revive
func (p *ProvidersService) GetProviderProjections(ctx context.Context, uid string, filter v1.ProviderFilter) (map[v1.Dim][]v1.Embedding, v1.Page, error) {
	p.db.RLock()
	defer p.db.RUnlock()
	count := 0
	if p.db.Closed {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	provider, ok := p.db.store[uid]
	if !ok {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
	}
	offset, ok := filter.Offset.(int)
	if !ok {
		offset = 0
	}
	if dim := filter.Dim; dim != nil {
		if *dim != v1.Dim2D && *dim != v1.Dim3D {
			return nil, v1.Page{Count: &count},
				v1.Errorf(v1.EINVALID, "invalid dimension %v for provider %q", *dim, uid)
		}

		projStore := provider[proj].(map[v1.Dim][]v1.Embedding)
		newProjections := getDimProjections(projStore, *dim)
		count = len(newProjections)

		return map[v1.Dim][]v1.Embedding{
			*filter.Dim: paging.ApplyOffsetLimit(newProjections, offset, filter.Limit).([]v1.Embedding),
		}, v1.Page{Count: &count}, nil
	}

	projections := provider[proj].(map[v1.Dim][]v1.Embedding)
	newProjections2D := getDimProjections(projections, v1.Dim2D)
	newProjections3D := getDimProjections(projections, v1.Dim3D)
	count = len(newProjections2D)

	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: paging.ApplyOffsetLimit(newProjections2D, offset, filter.Limit).([]v1.Embedding),
		v1.Dim3D: paging.ApplyOffsetLimit(newProjections3D, offset, filter.Limit).([]v1.Embedding),
	}, v1.Page{Count: &count}, nil
}

// UpdateProviderEmbeddings updates embeddings of a specific provider.
// nolint:revive
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, embeds []v1.Embedding, prj v1.Projection) ([]v1.Embedding, error) {
	p.db.Lock()
	defer p.db.Unlock()
	if p.db.Closed {
		return nil, ErrDBClosed
	}

	provider, ok := p.db.store[uid]
	if !ok {
		return nil, v1.Errorf(v1.ENOTFOUND, "provider %s not found", uid)
	}
	embs := provider[emb].([]v1.Embedding)
	newEmbs := make([]v1.Embedding, len(embs))
	copy(newEmbs, embs)
	newEmbs = append(newEmbs, embeds...)

	prjs, err := projection.Compute(newEmbs, prj)
	if err != nil {
		return nil, err
	}

	provider[emb] = newEmbs
	provider[proj] = prjs

	return embeds, nil
}

// DropProviderEmbeddings drops all embeddings for the provider with the given uid.
// NOTE: this obviously also drops the projections, as keeping them would make no sense
// since there would be no embeddings to associate them with.
// nolint:revive
func (p *ProvidersService) DropProviderEmbeddings(ctx context.Context, uid string) error {
	p.db.Lock()
	defer p.db.Unlock()
	if p.db.Closed {
		return v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	provider, ok := p.db.store[uid]
	if !ok {
		return v1.Errorf(v1.ENOTFOUND, "provider %q not found", uid)
	}
	provider[emb] = []v1.Embedding{}
	provider[proj] = map[v1.Dim][]v1.Embedding{}
	return nil
}

// ComputeProviderProjections recomputes all projections from scratch for the provider with the given UID.
// nolint:revive
func (p *ProvidersService) ComputeProviderProjections(ctx context.Context, uid string, prj v1.Projection) error {
	p.db.Lock()
	defer p.db.Unlock()
	if p.db.Closed {
		return v1.Errorf(v1.EINTERNAL, "%v", ErrDBClosed)
	}

	provider, ok := p.db.store[uid]
	if !ok {
		return v1.Errorf(v1.ENOTFOUND, "provider %s not found", uid)
	}
	embs := provider[emb].([]v1.Embedding)

	prjs, err := projection.Compute(embs, prj)
	if err != nil {
		return err
	}

	provider[proj] = prjs

	return nil
}

func getDimProjections(projections map[v1.Dim][]v1.Embedding, dim v1.Dim) []v1.Embedding {
	newProjections := make([]v1.Embedding, len(projections[dim]))
	copy(newProjections, projections[dim])
	return newProjections
}
