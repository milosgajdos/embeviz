package internal

import (
	"fmt"
	"strings"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/go-embeddings/document/text"
)

func getIndices(chunk, s string, startIdx int) []int {
	index := strings.Index(s[startIdx:], chunk)
	if index == -1 {
		return []int{}
	}
	return []int{index + startIdx, index + startIdx + len(chunk)}
}

// GetChunksIndices returns indices of each chunk in chunks inside s.
func GetChunksIndices(chunks []string, s string) [][]int {
	startIdx := 0
	result := make([][]int, 0, len(chunks))
	for _, chunk := range chunks {
		indices := getIndices(chunk, s, startIdx)
		if len(indices) > 0 {
			startIdx = indices[0]
			result = append(result, indices)
		}
	}
	return result
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

	return GetChunksIndices(rs.Split(req.Input), req.Input), nil
}
