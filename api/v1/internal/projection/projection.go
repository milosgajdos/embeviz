package projection

import (
	"errors"
	"fmt"

	"github.com/danaugrs/go-tsne/tsne"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"golang.org/x/exp/maps"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

var projDimToNum = map[v1.Dim]int{
	v1.Dim2D: 2,
	v1.Dim3D: 3,
}

// PCA computes PCA vectors and projects the original embeddings to the given dimension.
// It returns a new slice of embeddings of the same size as the origin al embeddings,
// but with the given dimension dropped to the given dimension.
func PCA(embs []v1.Embedding, projDim v1.Dim) ([]v1.Embedding, error) {
	dim := projDimToNum[projDim]
	pcas := make([]v1.Embedding, 0, len(embs))

	if embDim := len(embs[0].Values); embDim <= dim {
		return nil, fmt.Errorf("insufficient embedding dimension: %d, needs at least: %d", embDim, dim+1)
	}

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
		metadata := map[string]any{}
		if embs[i].Metadata != nil {
			metadata = maps.Clone(embs[i].Metadata)
		}
		metadata["projection"] = v1.PCA
		pcas = append(pcas, v1.Embedding{
			Values:   proj.RawRowView(i),
			Metadata: metadata,
		})
	}

	return pcas, nil
}

// TSNE calculates tsne projecion of the given embeddings into the given dimension.
// It returns a new slice of embeddings of the same size as the origin al embeddings,
// but with the given dimension dropped to the given dimension.
func TSNE(embs []v1.Embedding, projDim v1.Dim) ([]v1.Embedding, error) {
	dim := projDimToNum[projDim]
	tsnes := make([]v1.Embedding, 0, len(embs))

	if embDim := len(embs[0].Values); embDim <= dim {
		return nil, fmt.Errorf("insufficient embedding dimension: %d, needs at least: %d", embDim, dim+1)
	}

	// NOTE: these are somewhat randomly picked hyperparams YMMV
	perplexity, learningRate := float64(300), float64(300)
	if dim == 3 {
		perplexity, learningRate = float64(500), float64(500)
	}

	mx := mat.NewDense(len(embs), len(embs[0].Values), nil)
	for i, e := range embs {
		mx.SetRow(i, e.Values)
	}

	t := tsne.NewTSNE(dim, perplexity, learningRate, 300, true)
	resMat := t.EmbedData(mx, nil)
	d := mat.DenseCopyOf(resMat)

	for i := range embs {
		metadata := map[string]any{}
		if embs[i].Metadata != nil {
			metadata = maps.Clone(embs[i].Metadata)
		}
		metadata["projection"] = v1.TSNE
		tsnes = append(tsnes, v1.Embedding{
			Values:   d.RawRowView(i),
			Metadata: metadata,
		})
	}

	return tsnes, nil
}

// Compute computes p projections (2D and 3D) for embeddings embs and returns them.
func Compute(embs []v1.Embedding, p v1.Projection) (map[v1.Dim][]v1.Embedding, error) {
	if len(embs) == 0 {
		return map[v1.Dim][]v1.Embedding{
			v1.Dim2D: {},
			v1.Dim3D: {},
		}, nil
	}
	var (
		err    error
		proj2D []v1.Embedding
		proj3D []v1.Embedding
	)
	// Calculate projection
	switch p {
	case v1.PCA:
		proj2D, err = PCA(embs, v1.Dim2D)
		if err != nil {
			return nil, err
		}
		proj3D, err = PCA(embs, v1.Dim3D)
		if err != nil {
			return nil, err
		}
	case v1.TSNE:
		proj2D, err = TSNE(embs, v1.Dim2D)
		if err != nil {
			return nil, err
		}
		proj3D, err = TSNE(embs, v1.Dim3D)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid projection: %v", p)

	}
	return map[v1.Dim][]v1.Embedding{
		v1.Dim2D: proj2D,
		v1.Dim3D: proj3D,
	}, nil
}
