package qdrant

import (
	"context"
	"errors"

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

var (
	ErrMissingVectorSize     = errors.New("ErrMissingVectorSize")
	ErrInvalidVectorSize     = errors.New("ErrInvalidVectorSize")
	ErrInvalidVectorDistance = errors.New("ErrInvalidVectorDistance")
)

// ProvidersService allows to store data in qdrant vector store.
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
// It creates a new qdrant collection and raturns the new provider.
// The collection name is the same as the UUID of the provider.
func (p *ProvidersService) AddProvider(ctx context.Context, name string, md map[string]any) (*v1.Provider, error) {
	size, ok := md["size"]
	if !ok {
		return nil, v1.Errorf(v1.EINVALID, "%v", ErrMissingVectorSize)
	}
	vectorSize, ok := size.(uint64)
	if !ok {
		return nil, v1.Errorf(v1.EINVALID, "%v: %v", ErrInvalidVectorSize, vectorSize)
	}

	var vectorDistance pb.Distance
	dist, ok := md["distance"]
	if !ok {
		vectorDistance = defaultDistance
	} else {
		vectorDistance, ok = dist.(pb.Distance)
		if !ok {
			return nil, v1.Errorf(v1.EINVALID, "%v: %v", ErrInvalidVectorDistance, dist)
		}
	}

	resp, err := p.db.col.ListAliases(ctx, &pb.ListAliasesRequest{})
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "ListAliases error: %v", err)
	}

	for _, a := range resp.Aliases {
		if a.AliasName == name {
			md, err := p.getProviderMetadata(ctx, a.CollectionName)
			if err != nil {
				return nil, v1.Errorf(v1.EINTERNAL, "GetCollectionInfo error: %v", err)
			}
			return &v1.Provider{
				UID:      a.CollectionName,
				Name:     name,
				Metadata: md,
			}, nil
		}
	}

	uid := uuid.New().String()
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	_, err = p.db.col.Create(ctx, &pb.CreateCollection{
		CollectionName: uid,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_ParamsMap{
				ParamsMap: &pb.VectorParamsMap{
					Map: map[string]*pb.VectorParams{
						// NOTE(milosgajdos): empty name vector
						// is the "default" point vector.
						"": {
							Size:     vectorSize,
							Distance: vectorDistance,
						},
						"2D": {
							Size:     2,
							Distance: vectorDistance,
						},
						"3D": {
							Size:     3,
							Distance: vectorDistance,
						},
					},
				},
			},
		},
		OptimizersConfig: &pb.OptimizersConfigDiff{
			// https://qdrant.tech/documentation/concepts/optimizer/#merge-optimizer
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "CreateCollection error %v", err)
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
		return nil, v1.Errorf(v1.EINTERNAL, "UpdateAliases error %v", err)
	}

	return &v1.Provider{
		UID:      uid,
		Name:     name,
		Metadata: md,
	}, nil
}

// GetProviders returns a list of providers filtered by filter.
// NOTE: this does not populate metadata in v1.Provider because qdrant does not allow
// storing any metadata about the collections; only about the data stored in collections.
func (p *ProvidersService) GetProviders(ctx context.Context, filter v1.ProviderFilter) ([]*v1.Provider, v1.Page, error) {
	count := 0
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	resp, err := p.db.col.ListAliases(ctx, &pb.ListAliasesRequest{})
	if err != nil {
		return nil, v1.Page{Count: &count}, v1.Errorf(v1.EINTERNAL, "ListAliases error %v", err)
	}

	providers := map[string]string{}

	for _, a := range resp.Aliases {
		if _, ok := providers[a.CollectionName]; ok {
			continue
		}
		providers[a.CollectionName] = a.AliasName
	}

	px := make([]*v1.Provider, 0, len(providers))
	for uid, alias := range providers {
		md, err := p.getProviderMetadata(ctx, uid)
		if err != nil {
			return nil, v1.Page{Count: &count}, v1.Errorf(v1.EINTERNAL, "GetCollectionInfo error: %v", err)
		}
		px = append(px, &v1.Provider{
			UID:      uid,
			Name:     alias,
			Metadata: md,
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
// NOTE: this does not populate metadata in v1.Provider because qdrant does not allow
// storing any metadata about the collections; only about the data stored in collections.
func (p *ProvidersService) GetProviderByUID(ctx context.Context, uid string) (*v1.Provider, error) {
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	// fetch aliases for the given collection
	resp, err := p.db.httpClient.AliasList(ctx, uid)
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "AliasList error %v", err)
	}

	var alias string
	if len(resp.Result.Aliases) == 0 {
		return nil, v1.Errorf(v1.ENOTFOUND, "provider %s not found", uid)
	}
	alias = resp.Result.Aliases[0].Name

	md, err := p.getProviderMetadata(ctx, uid)
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "GetCollectionInfo error: %v", err)
	}

	return &v1.Provider{
		UID:      uid,
		Name:     alias,
		Metadata: md,
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

	page := v1.Page{}

	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, page, v1.Errorf(v1.EINTERNAL, "Scroll error %v", err)
	}

	points := resp.GetResult()
	embs := make([]v1.Embedding, len(points))

	for _, p := range points {
		// NOTE: this is a bit counter-intuitive; point is a default
		// vector which is essentially an unnamed vector in the vector map
		vec := p.GetVectors().GetVectors()
		if vec != nil {
			vals := make([]float64, 0, len(vec.Vectors[""].Data))
			for _, val := range vec.Vectors[""].Data {
				vals = append(vals, float64(val))
			}
			embs = append(embs, v1.Embedding{
				UID:      p.Id.GetUuid(),
				Values:   vals,
				Metadata: payload2Meta(p.GetPayload()),
			})
		}
	}

	if resp.NextPageOffset != nil {
		next := resp.NextPageOffset.GetUuid()
		page.Next = &next
	}

	return embs, page, nil
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

	page := v1.Page{}

	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	resp, err := p.db.pts.Scroll(ctx, req)
	if err != nil {
		return nil, page, v1.Errorf(v1.EINTERNAL, "Scroll error %v", err)
	}

	points := resp.GetResult()

	if resp.NextPageOffset != nil {
		next := resp.NextPageOffset.GetUuid()
		page.Next = &next
	}

	if dim := filter.Dim; dim != nil {
		if *dim != v1.Dim2D && *dim != v1.Dim3D {
			return nil, page, v1.Errorf(v1.EINVALID, "invalid dimension %v for provider %q", *dim, uid)
		}
		projs := make([]v1.Embedding, len(points))

		for _, p := range points {
			// NOTE: we call GetVectors twice because we use
			// NamedVectors so we need to dig in 2 levels down.
			vecs := p.GetVectors().GetVectors()
			if vecs != nil {
				vals := getVecVals(vecs, string(*dim))
				projs = append(projs, v1.Embedding{
					UID:      p.Id.GetUuid(),
					Values:   vals,
					Metadata: payload2Meta(p.GetPayload()),
				})
			}
		}
		return map[v1.Dim][]v1.Embedding{*filter.Dim: projs}, page, nil
	}

	res2DProjs := make([]v1.Embedding, 0, len(points))
	res3DProjs := make([]v1.Embedding, 0, len(points))

	for _, p := range points {
		// NOTE: we call GetVectors twice because we use
		// NamedVectors so we need to dig in 2 levels down.
		vecs := p.GetVectors().GetVectors()
		// skip if no named vectors exist for this point
		// this means there are no projections for this embedding.
		if vecs != nil {
			// 2D projections
			proj2DVals := getVecVals(vecs, string(v1.Dim2D))
			res2DProjs = append(res2DProjs, v1.Embedding{
				UID:      p.Id.GetUuid(),
				Values:   proj2DVals,
				Metadata: payload2Meta(p.GetPayload()),
			})
			// 3D projections
			proj3DVals := getVecVals(vecs, string(v1.Dim3D))
			res3DProjs = append(res3DProjs, v1.Embedding{
				UID:      p.Id.GetUuid(),
				Values:   proj3DVals,
				Metadata: payload2Meta(p.GetPayload()),
			})
		}
	}
	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: res2DProjs,
		v1.Dim3D: res3DProjs,
	}, page, nil
}

// UpdateProviderEmbeddings generates embeddings for the provider with the given uid.
func (p *ProvidersService) UpdateProviderEmbeddings(ctx context.Context, uid string, embeds []v1.Embedding, proj v1.Projection) ([]v1.Embedding, error) {
	// fetch all points so we can compute projections
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)

	upsertPoints := make([]*pb.PointStruct, 0, len(embeds))
	for _, e := range embeds {
		// create a new embedding point
		data := make([]float32, 0, len(e.Values))
		for _, val := range e.Values {
			data = append(data, float32(val))
		}
		pointUID := e.UID
		if pointUID == "" {
			pointUID = uuid.NewString()
		}
		md := make(map[string]*pb.Value)
		for k, v := range e.Metadata {
			// FIXME: v can be of any type
			if stringVal, ok := v.(string); ok {
				md[k] = &pb.Value{
					Kind: &pb.Value_StringValue{
						StringValue: stringVal,
					},
				}
			}
		}
		upsertPoints = append(upsertPoints, &pb.PointStruct{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{
					Uuid: pointUID,
				},
			},
			Vectors: &pb.Vectors{
				VectorsOptions: &pb.Vectors_Vectors{
					Vectors: &pb.NamedVectors{
						Vectors: map[string]*pb.Vector{
							"":   {Data: data},
							"2D": {Data: []float32{0, 0}},
							"3D": {Data: []float32{0, 0, 0}},
						},
					},
				},
			},
			Payload: md,
		})
	}

	// TODO: can we avoid upserting duplicate points ?
	// duplicate meaning similar vectors or the same vectors i.e.
	// that have the same data not necessarily the same UUID
	waitUpsert := true
	if _, err := p.db.pts.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: uid,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	}); err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "Upsert error %v", err)
	}

	// Collect all embeddings and calculate projections.
	// NOTE: tread carefully, as this can shit memory pants on large collections!
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

	embs := []v1.Embedding{}

	for {
		resp, err := p.db.pts.Scroll(ctx, req)
		if err != nil {
			return nil, v1.Errorf(v1.EINTERNAL, "Scroll error %v", err)
		}
		next := resp.NextPageOffset

		for _, p := range resp.GetResult() {
			vecs := p.GetVectors().GetVectors()
			if vecs != nil {
				vals := make([]float64, 0, len(vecs.Vectors[""].Data))
				for _, val := range vecs.Vectors[""].Data {
					vals = append(vals, float64(val))
				}
				embs = append(embs, v1.Embedding{
					UID:      p.Id.GetUuid(),
					Values:   vals,
					Metadata: payload2Meta(p.GetPayload()),
				})
			}
		}
		// stop paging we're done
		if next == nil {
			break
		}
		req.Offset = next
	}

	projs, err := projection.Compute(embs, proj)
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "Compute error %v", err)
	}

	pointVecs := make([]*pb.PointVectors, 0, len(embs))

	for i, emb := range embs {
		pv := &pb.PointVectors{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{
					Uuid: emb.UID,
				},
			},
		}
		namedVecs := make(map[string]*pb.Vector)
		for dim, dimProjs := range projs {
			data := make([]float32, 0, len(dimProjs[i].Values))
			for _, val := range dimProjs[i].Values {
				data = append(data, float32(val))
			}
			namedVecs[string(dim)] = &pb.Vector{
				Data: data,
			}
		}
		pv.Vectors = &pb.Vectors{
			VectorsOptions: &pb.Vectors_Vectors{
				Vectors: &pb.NamedVectors{
					Vectors: namedVecs,
				},
			},
		}
		pointVecs = append(pointVecs, pv)
	}

	if _, err := p.db.pts.UpdateVectors(ctx, &pb.UpdatePointVectors{
		CollectionName: uid,
		Wait:           &waitUpsert,
		Points:         pointVecs,
	}); err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "UpdateVectors error %v", err)
	}

	return embeds, nil
}

// DropProviderEmbeddings drops all provider embeddings from the store
func (p *ProvidersService) DropProviderEmbeddings(ctx context.Context, uid string) error {
	ctx = metadata.NewOutgoingContext(ctx, p.db.md)
	// * retrieven collection
	// * grab the vector config
	col, err := p.db.col.Get(ctx, &pb.GetCollectionInfoRequest{CollectionName: uid})
	if err != nil {
		return v1.Errorf(v1.EINTERNAL, "GetCollection error: %v", err)
	}
	vecConfig := col.Result.Config.Params.GetVectorsConfig()

	// * fetch aliases
	resp, err := p.db.httpClient.AliasList(ctx, uid)
	if err != nil {
		return v1.Errorf(v1.EINTERNAL, "AliasList error: %v", err)
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
		return v1.Errorf(v1.EINTERNAL, "DeleteAliases error: %v", err)
	}

	// * drop collection
	if _, err := p.db.col.Delete(ctx, &pb.DeleteCollection{
		CollectionName: uid,
	}); err != nil {
		return v1.Errorf(v1.EINTERNAL, "DeleteCollection: %v", err)
	}

	// create a new collection
	if _, err = p.db.col.Create(ctx, &pb.CreateCollection{
		CollectionName: uid,
		VectorsConfig:  vecConfig,
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	}); err != nil {
		return v1.Errorf(v1.EINTERNAL, "CreateCollection: %v", err)
	}

	if _, err := p.db.col.UpdateAliases(ctx, &pb.ChangeAliases{Actions: createAliases}); err != nil {
		return v1.Errorf(v1.EINTERNAL, "UpdateAliases error: %v", err)
	}

	return nil
}

// ComputeProviderProjections recomputes all projections from scratch for the provider with the given UID.
func (p *ProvidersService) ComputeProviderProjections(ctx context.Context, uid string, proj v1.Projection) error {
	// NOTE: tread carefully, as this can shit memory pants on large collections!
	// fetch all points so we can compute projections
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

	embs := []v1.Embedding{}

	for {
		resp, err := p.db.pts.Scroll(ctx, req)
		if err != nil {
			return v1.Errorf(v1.EINTERNAL, "Scroll error %v", err)
		}
		next := resp.NextPageOffset

		for _, p := range resp.GetResult() {
			vecs := p.GetVectors().GetVectors()
			if vecs != nil {
				vals := make([]float64, 0, len(vecs.Vectors[""].Data))
				for _, val := range vecs.Vectors[""].Data {
					vals = append(vals, float64(val))
				}
				embs = append(embs, v1.Embedding{
					UID:      p.Id.GetUuid(),
					Values:   vals,
					Metadata: payload2Meta(p.GetPayload()),
				})
			}
		}
		// stop paging we're done
		if next == nil {
			break
		}
		req.Offset = next
	}

	projs, err := projection.Compute(embs, proj)
	if err != nil {
		return v1.Errorf(v1.EINTERNAL, "Compute error %v", err)
	}

	pointVecs := make([]*pb.PointVectors, 0, len(embs))

	for i, emb := range embs {
		pv := &pb.PointVectors{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{
					Uuid: emb.UID,
				},
			},
		}
		namedVecs := make(map[string]*pb.Vector)
		for dim, dimProjs := range projs {
			data := make([]float32, 0, len(dimProjs[i].Values))
			for _, val := range dimProjs[i].Values {
				data = append(data, float32(val))
			}
			namedVecs[string(dim)] = &pb.Vector{
				Data: data,
			}
		}
		pv.Vectors = &pb.Vectors{
			VectorsOptions: &pb.Vectors_Vectors{
				Vectors: &pb.NamedVectors{
					Vectors: namedVecs,
				},
			},
		}
		pointVecs = append(pointVecs, pv)
	}

	waitUpsert := true
	if _, err := p.db.pts.UpdateVectors(ctx, &pb.UpdatePointVectors{
		CollectionName: uid,
		Wait:           &waitUpsert,
		Points:         pointVecs,
	}); err != nil {
		return v1.Errorf(v1.EINTERNAL, "UpdateVectors error %v", err)
	}

	return nil
}

// getProviderMetadata returns metadata for the provider with the given uid.
func (p *ProvidersService) getProviderMetadata(ctx context.Context, uid string) (map[string]any, error) {
	col, err := p.db.col.Get(ctx, &pb.GetCollectionInfoRequest{CollectionName: uid})
	if err != nil {
		return nil, v1.Errorf(v1.EINTERNAL, "GetCollection error: %v", err)
	}
	params := col.Result.Config.Params.GetVectorsConfig().GetParamsMap()

	md := make(map[string]any)

	md["size"] = params.Map[""].Size
	md["distance"] = params.Map[""].Distance

	return md, nil
}

func payload2Meta(payload map[string]*pb.Value) map[string]any {
	md := make(map[string]any)
	for k, v := range payload {
		switch v := v.GetKind().(type) {
		case *pb.Value_DoubleValue:
			md[k] = v.DoubleValue
		case *pb.Value_IntegerValue:
			md[k] = v.IntegerValue
		case *pb.Value_NullValue:
			md[k] = v.NullValue
		case *pb.Value_BoolValue:
			md[k] = v.BoolValue
		case *pb.Value_StructValue:
			md[k] = v.StructValue
		case *pb.Value_StringValue:
			md[k] = v.StringValue
		default:
			continue
		}
	}
	return md
}
