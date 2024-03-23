package internal

import (
	"fmt"
	"strings"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/go-embeddings/document/text"
)

// GetChunkIndices returns indices of each chunk in chunks inside s.
func GetChunkIndices(chunks []string, s string) [][]int {
	var indices [][]int //nolint:prealloc

	// Iterate over each chunk
	for _, chunk := range chunks {
		start := strings.Index(s, chunk)
		if start == -1 {
			continue // Chunk not found in the string, skip
		}
		end := start + len(chunk)
		indices = append(indices, []int{start, end})
	}

	return indices
}

// GetChunks chunks the input and returns chunk indices of each chunk.
// If the req is nil it rerurns error.
func GetChunks(req *v1.ChunkingInput) ([][]int, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: %v", req)
	}

	if len(req.Input) == 0 {
		return [][]int{}, nil
	}

	s := text.NewSplitter().
		WithChunkSize(req.Options.Size).
		WithChunkOverlap(req.Options.Overlap).
		WithTrimSpace(req.Options.Trim).
		WithKeepSep(req.Options.Sep)

	rs := text.NewRecursiveCharSplitter().
		WithSplitter(s)

	return GetChunkIndices(rs.Split(req.Input), req.Input), nil
}
