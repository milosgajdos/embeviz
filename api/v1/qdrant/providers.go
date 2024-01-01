package qdrant

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/internal/paging"
	"github.com/milosgajdos/embeviz/api/v1/internal/projection"
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

	createAliases := []*pb.AliasOperations{
		{
			Action: &pb.AliasOperations_CreateAlias{
				CreateAlias: &pb.CreateAlias{
					CollectionName: uid,
					AliasName:      name,
				},
			},
		},
	}

	if _, err := p.db.col.UpdateAliases(ctx, &pb.ChangeAliases{Actions: createAliases}); err != nil {
		return nil, err
	}

	return &v1.Provider{
		UID:      uid,
		Name:     name,
		Metadata: md,
	}, nil
}

// GetProviders returns a list of providers filtered by filter.
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, v1.Page, error) {
	count := 0
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	resp, err := p.db.col.ListAliases(ctx, &pb.ListAliasesRequest{})
	if err != nil {
		return nil, v1.Page{Count: &count}, err
	}

	providers := map[string]string{}

	for _, a := range resp.Aliases {
		fmt.Println("Collection", a.CollectionName)
		if _, ok := providers[a.CollectionName]; ok {
			continue
		}
		providers[a.CollectionName] = a.AliasName
	}

	px := make([]*v1.Provider, 0, len(providers))
	for uid, alias := range providers {
		px = append(px, &v1.Provider{
			UID:  uid,
			Name: alias,
		})
	}
	count = len(px)

	offset, ok := filter.Offset.(int)
	if !ok {
		offset = 0
	}

	return paging.ApplyOffsetLimit(px, offset, filter.Limit).([]*v1.Provider), v1.Page{Count: &count}, nil
}

// GetProviderByUID returns the provider with the given uid.
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	// * fetch aliases for the given collection
	resp, err := p.db.httpClient.AliasList(ctx, uid)
	if err != nil {
		return nil, err
	}

	var alias string
	if len(resp.Result.Aliases) == 0 {
		return nil, v1.Errorf(v1.ENOTFOUND, "provider %s not found", uid)
	}
	alias = resp.Result.Aliases[0].Name

	return &v1.Provider{
		UID:  uid,
		Name: alias,
	}, nil
}

// GetProviderEmbeddings returns embeddings for the provider with the given uid.
func (p *ProvidersService) GetProviderEmbeddings(ctx context.Context, uid string, filter v1.ProviderFilter) ([]v1.Embedding, v1.Page, error) {
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

	count := 0
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, v1.Page{Count: &count}, err
	}
	next := resp.NextPageOffset.String()

	points := resp.GetResult()
	embs := make([]v1.Embedding, len(points))

	for _, p := range points {
		vec := p.GetVectors().GetVector()
		// TODO: grab metadata from p.Payload
		vals := make([]float64, 0, len(vec.Data))
		for _, val := range vec.Data {
			vals = append(vals, float64(val))
		}
		embs = append(embs, v1.Embedding{
			UID:    p.Id.String(),
			Values: vals,
		})
	}

	return embs, v1.Page{Next: &next}, nil
}

// GetProviderProjections returns embeddings projections for the provider with the given uid.
func (p *ProvidersService) GetProviderProjections(ctx context.Context, uid string, filter v1.ProviderFilter) (map[v1.Dim][]v1.Embedding, v1.Page, error) {
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

	count := 0
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, v1.Page{Count: &count}, err
	}
	next := resp.NextPageOffset.String()

	points := resp.GetResult()

	if dim := filter.Dim; dim != nil {
		if *dim != v1.Dim2D && *dim != v1.Dim3D {
			return nil, v1.Page{Count: &count},
				v1.Errorf(v1.EINVALID, "invalid dimension %v for provider %q", *dim, uid)
		}
		projs := make([]v1.Embedding, len(points))

		for _, p := range points {
			// NOTE: we call GetVectors twice because we use
			// NamedVectors so we need to dig in 2 levels down.
			vecs := p.GetVectors().GetVectors()
			if vecs != nil {
				vals := getNamedVecVals(vecs, string(*dim))
				projs = append(projs, v1.Embedding{
					Values: vals,
				})
			}
		}
		return map[v1.Dim][]v1.Embedding{*filter.Dim: projs}, v1.Page{Next: &next}, nil
	}

	res2DProjs := make([]v1.Embedding, len(points))
	res3DProjs := make([]v1.Embedding, len(points))

	for _, p := range points {
		// NOTE: we call GetVectors twice because we use
		// NamedVectors so we need to dig in 2 levels down.
		vecs := p.GetVectors().GetVectors()
		// skip if no named vectors exist for this point
		// this means there are no projections for this embedding.
		if vecs != nil {
			// 2D projections
			proj2DVals := getNamedVecVals(vecs, string(v1.Dim2D))
			res2DProjs = append(res2DProjs, v1.Embedding{
				Values: proj2DVals,
			})
			// 3D projections
			proj3DVals := getNamedVecVals(vecs, string(v1.Dim3D))
			res3DProjs = append(res3DProjs, v1.Embedding{
				Values: proj3DVals,
			})
		}
	}
	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: res2DProjs,
		v1.Dim3D: res3DProjs,
	}, v1.Page{Next: &next}, nil
}

// UpdateProviderEmbeddings generates embeddings for the provider with the given uid.
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, embed v1.Embedding, proj v1.Projection) (*v1.Embedding, error) {
	// fetch all points so we can compute PCA
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
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

	// create a new embedding point
	data := make([]float32, 0, len(embed.Values))
	for _, val := range embed.Values {
		data = append(data, float32(val))
	}
	upsertPoints := []*pb.PointStruct{
		{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{
					Uuid: uuid.NewString(),
				},
			},
			Vectors: &pb.Vectors{
				VectorsOptions: &pb.Vectors_Vector{
					Vector: &pb.Vector{
						Data: data,
					},
				},
			},
			Payload: map[string]*pb.Value{},
		},
	}
	waitUpsert := true
	if _, err := p.db.pts.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: uid,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	}); err != nil {
		return nil, err
	}

	// Collect all embeddings and calculate projections.
	embs := []v1.Embedding{}
	for {
		resp, err := p.db.pts.Scroll(ctx, req)
		if err != nil {
			return nil, err
		}
		next := resp.NextPageOffset

		for _, p := range resp.GetResult() {
			vec := p.GetVectors().GetVector()
			// TODO: grab metadata from p.Payload
			vals := make([]float64, 0, len(vec.Data))
			for _, val := range vec.Data {
				vals = append(vals, float64(val))
			}
			embs = append(embs, v1.Embedding{
				UID:    p.Id.String(),
				Values: vals,
			})
		}
		// stop paging we're done
		if next == nil {
			break
		}
		req.Offset = next
	}

	projs, err := projection.Compute(embs, proj)
	if err != nil {
		return nil, err
	}

	points := make([]*pb.PointStruct, 0, len(projs))

	for dim, dimProjs := range projs {
		for i := range dimProjs {
			data := make([]float32, 0, len(dimProjs[i].Values))
			for _, val := range dimProjs[i].Values {
				data = append(data, float32(val))
			}
			points = append(points, &pb.PointStruct{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{
						Uuid: embs[i].UID,
					},
				},
				Vectors: &pb.Vectors{
					VectorsOptions: &pb.Vectors_Vectors{
						Vectors: &pb.NamedVectors{
							Vectors: map[string]*pb.Vector{
								string(dim): {Data: data},
							},
						},
					},
				},
			})
		}
	}

	// NOTE: this might need to be replaced with
	// UpdateVectors, but we would als need to collect
	// projections into pb.PointVectors
	if _, err := p.db.pts.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: uid,
		Wait:           &waitUpsert,
		Points:         points,
	}); err != nil {
		return nil, err
	}

	return &embed, nil
}

// DropProviderEmbeddings drops all provider embeddings from the store
// nolint:revive
func (p *ProvidersService) DropProviderEmbeddings(ctx context.Context, uid string) error {
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	// * retrieven collection
	// * grab the vector config
	col, err := p.db.col.Get(ctx, &pb.GetCollectionInfoRequest{CollectionName: uid})
	if err != nil {
		return err
	}
	vecConfig := col.Result.Config.Params.GetVectorsConfig()

	// * fetch aliases
	resp, err := p.db.httpClient.AliasList(ctx, uid)
	if err != nil {
		return err
	}

	// actions for deleting aliases
	deleteAliases := make([]*pb.AliasOperations, 0, len(resp.Result.Aliases))
	// actions for re-creating aliases
	createAliases := make([]*pb.AliasOperations, 0, len(resp.Result.Aliases))

	for _, alias := range resp.Result.Aliases {
		deleteAliases = append(deleteAliases, &pb.AliasOperations{
			Action: &pb.AliasOperations_DeleteAlias{
				DeleteAlias: &pb.DeleteAlias{AliasName: alias.Name},
			},
		})
		createAliases = append(createAliases, &pb.AliasOperations{
			Action: &pb.AliasOperations_CreateAlias{
				CreateAlias: &pb.CreateAlias{
					CollectionName: uid,
					AliasName:      alias.Name,
				},
			},
		})
	}
	if _, err := p.db.col.UpdateAliases(ctx, &pb.ChangeAliases{Actions: deleteAliases}); err != nil {
		return err
	}

	// * drop collection
	if _, err := p.db.col.Delete(ctx, &pb.DeleteCollection{
		CollectionName: uid,
	}); err != nil {
		return err
	}

	// * create new collection
	if _, err = p.db.col.Create(ctx, &pb.CreateCollection{
		CollectionName: uid,
		VectorsConfig:  vecConfig,
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	}); err != nil {
		return err
	}

	if _, err := p.db.col.UpdateAliases(ctx, &pb.ChangeAliases{Actions: createAliases}); err != nil {
		return err
	}

	return nil
}

// ComputeProviderProjections drops existing projections and recomputes anew.
// nolint:revive
func (p *ProvidersService) ComputeProviderProjections(ctx context.Context, uid string, proj v1.Projection) error {
	return v1.Errorf(v1.ENOTIMPLEMENTED, "ComputeProviderProjections")
}
