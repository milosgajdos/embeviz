package qdrant

import (
	pb "github.com/qdrant/go-client/qdrant"
)

// getVecVals returns data values for the given vector dimension.
func getVecVals(vecs *pb.NamedVectors, dim string) []float64 {
	// TODO: check if that key actually exists
	// and maybe return some flag or error.
	vecData := vecs.Vectors[dim].Data
	vals := make([]float64, 0, len(vecData))
	for _, val := range vecData {
		vals = append(vals, float64(val))
	}
	return vals
}
