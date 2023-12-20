package projection

import (
	"errors"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"golang.org/x/exp/maps"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func PCA(embs []v1.Embedding, projDim v1.Dim) ([]v1.Embedding, error) {
	dim := projDimToNum[projDim]
	pcas := make([]v1.Embedding, 0, len(embs))

	mx := mat.NewDense(len(embs), len(embs[0].Values), nil)
	for i, e := range embs {
		mx.SetRow(i, e.Values)
	}
	r, _ := mx.Dims()
	// Keep extending matrix until we have enough
	// data to compute the PCA
	for r < dim {
		// NOTE: if there is only one embedding, we simply add an identical
		// vector to the matrix and do a PCA on it.
		vals := make([]float64, len(embs[0].Values))
		copy(vals, embs[0].Values)
		mx = mx.Grow(1, 0).(*mat.Dense)
		mx.SetRow(1, vals)
		r, _ = mx.Dims()
	}
	var pc stat.PC
	ok := pc.PrincipalComponents(mx, nil)
	if !ok {
		return nil, errors.New("failed pca")
	}
	var proj mat.Dense
	var vec mat.Dense
	pc.VectorsTo(&vec)
	proj.Mul(mx, vec.Slice(0, len(embs[0].Values), 0, dim))

	for i := range embs {
		pcas = append(pcas, v1.Embedding{
			Values:   proj.RawRowView(i),
			Metadata: maps.Clone(embs[i].Metadata),
		})
	}

	return pcas, nil
}
