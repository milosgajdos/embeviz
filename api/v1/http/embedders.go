package http

import (
	"context"
	"errors"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/go-embeddings"
	"github.com/milosgajdos/go-embeddings/cohere"
	"github.com/milosgajdos/go-embeddings/openai"
	"github.com/milosgajdos/go-embeddings/vertexai"
)

// FetchEmbeddings fetches embeddings using the provided embedder.
// It returns the fetched embedding or fails with error.
func FetchEmbeddings(ctx context.Context, embedder any, req *v1.EmbeddingUpdate) (*v1.Embedding, error) {
	var (
		vals []float64
		embs []*embeddings.Embedding
		err  error
	)
	switch p := embedder.(type) {
	case *vertexai.Client:
		embReq := &vertexai.EmbeddingRequest{
			Instances: []vertexai.Instance{
				{
					Content:  req.Text,
					TaskType: vertexai.RetrQueryTask,
				},
			},
			Params: vertexai.Params{
				AutoTruncate: false,
			},
		}
		embs, err = p.Embed(ctx, embReq)
		if err != nil {
			return nil, err
		}
	case *openai.Client:
		embReq := &openai.EmbeddingRequest{
			Input:          req.Text,
			Model:          openai.TextAdaV2,
			EncodingFormat: openai.EncodingFloat,
		}
		embs, err = p.Embed(ctx, embReq)
		if err != nil {
			return nil, err
		}
	case *cohere.Client:
		embReq := &cohere.EmbeddingRequest{
			Texts:     []string{req.Text},
			Model:     cohere.EnglishV3,
			InputType: cohere.ClusteringInput,
			Truncate:  cohere.NoneTrunc,
		}
		embs, err = p.Embed(ctx, embReq)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported provider")
	}
	if len(embs) > 0 {
		vals = make([]float64, len(embs[0].Vector))
		copy(vals, embs[0].Vector)
	}
	return &v1.Embedding{
		Values:   vals,
		Metadata: req.Metadata,
	}, nil
}
