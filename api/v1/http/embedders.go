package http

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/go-embeddings"
	"github.com/milosgajdos/go-embeddings/cohere"
	"github.com/milosgajdos/go-embeddings/document/text"
	"github.com/milosgajdos/go-embeddings/openai"
	"github.com/milosgajdos/go-embeddings/vertexai"
)

// FetchEmbeddings fetches embeddings using the provided embedder.
// It returns the fetched embeddings or fails with error.
func FetchEmbeddings(ctx context.Context, embedder any, req *v1.EmbeddingsUpdate) ([]v1.Embedding, error) {
	var (
		vals []float64
		embs []*embeddings.Embedding
		err  error
	)

	if len(req.Text) == 0 {
		return []v1.Embedding{}, nil
	}

	chunks := []string{req.Text}

	// chunk input data if requested
	if req.Chunking != nil {
		s := text.NewSplitter().
			WithChunkSize(req.Chunking.Size).
			WithChunkOverlap(req.Chunking.Overlap).
			WithTrimSpace(req.Chunking.Trim).
			WithKeepSep(req.Chunking.Sep)

		rs := text.NewRecursiveCharSplitter().
			WithSplitter(s)

		chunks = rs.Split(req.Text)
	}

	// TODO: make sure the input vals for Embedding request
	// isn't exceeding max number of tokens e.g.
	// for OpenAI: https://platform.openai.com/docs/api-reference/embeddings/create

	results := make([]v1.Embedding, 0, len(chunks))

	switch p := embedder.(type) {
	case *vertexai.Client:
		instances := make([]vertexai.Instance, 0, len(chunks))
		for _, chunk := range chunks {
			instances = append(instances, vertexai.Instance{
				Content:  chunk,
				TaskType: vertexai.RetrQueryTask,
			})
		}
		embReq := &vertexai.EmbeddingRequest{
			Instances: instances,
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
			Input:          chunks,
			Model:          openai.TextAdaV2,
			EncodingFormat: openai.EncodingFloat,
		}
		embs, err = p.Embed(ctx, embReq)
		if err != nil {
			return nil, err
		}
	case *cohere.Client:
		embReq := &cohere.EmbeddingRequest{
			Texts:     chunks,
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

	for _, emb := range embs {
		vals = make([]float64, len(emb.Vector))
		copy(vals, emb.Vector)
		r := v1.Embedding{
			UID:      uuid.NewString(),
			Values:   vals,
			Metadata: req.Metadata,
		}
		results = append(results, r)
	}

	return results, nil
}
