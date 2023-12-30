package qdrant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

var (
	defaultSegmentNumber uint64 = 2
	defaultDistance             = pb.Distance_Dot
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
// It creates a new collection and raturns the new provider.
// The collection name is the same as the UUID of the provider.
func (p *ProvidersService) AddProvider(ctx context.Context, name string, md map[string]any) (*v1.Provider, error) {
	size, ok := md["size"]
	if !ok {
		return nil, errors.New("missing vector size")
	}
	vectorSize, ok := size.(uint64)
	if !ok {
		return nil, errors.New("invalid vector size")
	}

	var distance pb.Distance
	dist, ok := md["distance"]
	if !ok {
		distance = defaultDistance
	} else {
		distance, ok = dist.(pb.Distance)
		if !ok {
			return nil, errors.New("invalid vector distance")
		}
	}

	uid := uuid.New().String()
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	_, err := p.db.col.Create(ctx, &pb.CreateCollection{
		CollectionName: uid,
		VectorsConfig: &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{
			Params: &pb.VectorParams{
				Size:     vectorSize,
				Distance: distance,
			},
		}},
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		return nil, err
	}

	// TODO: create alias

	return &v1.Provider{
		UID:      uid,
		Name:     name,
		Metadata: md,
	}, nil
}

// GetProviders returns a list of providers filtered by filter.
// nolint:revive
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, int, error) {
	// TODO: requires collection alias
	return nil, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviders")
}

// GetProviderByUID returns the provider with the given uuid.
// nolint:revive
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	// TODO: requires collection alias
	return nil, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderByUID")
}

// GetProviderEmbeddings returns embeddings for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) GetProviderEmbeddings(ctx context.Context, uid string, filter v1.ProviderFilter) ([]v1.Embedding, int, error) {
	req := &pb.ScrollPoints{
		CollectionName: uid,
		WithVectors: &pb.WithVectorsSelector{
			SelectorOptions: &pb.WithVectorsSelector_Enable{
				Enable: true,
			},
		},
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	}
	// NOTE: qdrant doesn't do offset by numbers
	// but instead it offsets by Point IDs which
	// can be either int or string <insert_sad_emoji>
	if filter.Limit > 0 {
		limit := uint32(filter.Limit)
		req.Limit = &limit
	}

	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	points := resp.GetResult()
	embs := make([]v1.Embedding, len(points))

	for _, p := range points {
		// NOTE: we call GetVectors twice because we use
		// NamedVectors so we need to dig in 2 levels down.
		vecs := p.GetVectors().GetVectors()
		vecVals := vecs.Vectors["embs"].Data
		embVals := make([]float64, len(vecVals))
		for _, val := range vecVals {
			embVals = append(embVals, float64(val))
		}
		embs = append(embs, v1.Embedding{
			Values: embVals,
		})
	}

	return embs, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderEmbeddings")
}

// GetProviderProjections returns embeddings projections for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) GetProviderProjections(ctx context.Context, uid string, filter v1.ProviderFilter) (map[v1.Dim][]v1.Embedding, int, error) {
	req := &pb.ScrollPoints{
		CollectionName: uid,
		WithVectors: &pb.WithVectorsSelector{
			SelectorOptions: &pb.WithVectorsSelector_Enable{
				Enable: true,
			},
		},
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	}
	// NOTE: qdrant doesn't do offset by numbers
	// but instead it offsets by Point IDs which
	// can be either int or string <insert_sad_emoji>
	if filter.Limit > 0 {
		limit := uint32(filter.Limit)
		req.Limit = &limit
	}

	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	points := resp.GetResult()

	if dim := filter.Dim; dim != nil {
		if *dim != v1.Dim2D && *dim != v1.Dim3D {
			return nil, 0, v1.Errorf(v1.EINVALID, "invalid dimension %v for provider %q", *dim, uid)
		}
		projs := make([]v1.Embedding, len(points))

		for _, p := range points {
			// NOTE: we call GetVectors twice because we use
			// NamedVectors so we need to dig in 2 levels down.
			vecs := p.GetVectors().GetVectors()
			vals := getVecDimVals(vecs, *dim)
			projs = append(projs, v1.Embedding{
				Values: vals,
			})
		}
		return map[v1.Dim][]v1.Embedding{*filter.Dim: projs}, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderEmbeddings")
	}

	res2DProjs := make([]v1.Embedding, len(points))
	res3DProjs := make([]v1.Embedding, len(points))

	for _, p := range points {
		// NOTE: we call GetVectors twice because we use
		// NamedVectors so we need to dig in 2 levels down.
		vecs := p.GetVectors().GetVectors()

		proj2DVals := getVecDimVals(vecs, v1.Dim2D)
		res2DProjs = append(res2DProjs, v1.Embedding{
			Values: proj2DVals,
		})

		proj3DVals := getVecDimVals(vecs, v1.Dim3D)
		res3DProjs = append(res3DProjs, v1.Embedding{
			Values: proj3DVals,
		})
	}

	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: res2DProjs,
		v1.Dim3D: res3DProjs,
	}, 0, v1.Errorf(v1.ENOTIMPLEMENTED, "GetProviderEmbeddings")
}

// UpdateProviderEmbeddings generates embeddings for the provider with the given uid.
// nolint:revive
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, update v1.Embedding, proj v1.Projection) (*v1.Embedding, error) {
	// TODO: need to modify points https://qdrant.tech/documentation/concepts/points/#modify-points
	return nil, v1.Errorf(v1.ENOTIMPLEMENTED, "UpdateProviderEmbeddings")
}

// DropProviderEmbeddings drops all provider embeddings from the store
// nolint:revive
func (p *ProvidersService) DropProviderEmbeddings(ctx context.Context, uid string) error {
	// TODO:
	// * fetch an alias
	// * drop alias
	// * drop collection
	// * create new collection
	// * create alias
	return v1.Errorf(v1.ENOTIMPLEMENTED, "DropProviderEmbeddings")
}

// ComputeProviderProjections drops existing projections and recomputes anew.
// nolint:revive
func (p *ProvidersService) ComputeProviderProjections(ctx context.Context, uid string, proj v1.Projection) error {
	return v1.Errorf(v1.ENOTIMPLEMENTED, "ComputeProviderProjections")
}
