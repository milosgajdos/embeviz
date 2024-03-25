package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetChunkIndices(t *testing.T) {
	var testCases = []struct {
		input  string
		chunks []string
		exp    [][]int
	}{
		{"", []string{}, [][]int{}},
		{"foo", []string{}, [][]int{}},
		{"", []string{"foo", "bar"}, [][]int{}},
		{"こんにちは世界", []string{"こ", "世界"}, [][]int{{0, 3}, {15, 21}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("input=%q,chunks=%v", tc.input, tc.chunks), func(t *testing.T) {
			t.Parallel()
			indices := GetChunksIndices(tc.chunks, tc.input)
			if !reflect.DeepEqual(indices, tc.exp) {
				t.Errorf("expected: %v, got: %v", tc.exp, indices)
			}
		})
	}
}
