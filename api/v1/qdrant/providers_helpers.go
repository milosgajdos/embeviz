package qdrant

import (
	pb "github.com/qdrant/go-client/qdrant"
)

// getNamedVecVals returns data values for the given vector name.
func getNamedVecVals(namedVecs *pb.NamedVectors, name string) []float64 {
	// TODO: check if that key actually exists
	// and maybe return some flag or error.
	vecData := namedVecs.Vectors[name].Data
	vals := make([]float64, 0, len(vecData))
	for _, val := range vecData {
		vals = append(vals, float64(val))
	}
	return vals
}
