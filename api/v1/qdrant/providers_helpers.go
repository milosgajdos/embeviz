package qdrant

import (
	v1 "github.com/milosgajdos/embeviz/api/v1"
	pb "github.com/qdrant/go-client/qdrant"
)

// getVecDimVals returns data values for the given dimension from the named vectors.
func getVecDimVals(namedVecs *pb.NamedVectors, dim v1.Dim) []float64 {
	// TODO: check if that key actually exists
	// and maybe return some flag or error.
	vecData := namedVecs.Vectors[string(dim)].Data
	vals := make([]float64, len(vecData))
	for _, val := range vecData {
		vals = append(vals, float64(val))
	}
	return vals
}
