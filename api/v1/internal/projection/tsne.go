package projection

import (
	"github.com/danaugrs/go-tsne/tsne"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	"golang.org/x/exp/maps"
	"gonum.org/v1/gonum/mat"
)

func TSNE(embs []v1.Embedding, projDim v1.Dim) ([]v1.Embedding, error) {
	dim := projDimToNum[projDim]
	tsnes := make([]v1.Embedding, 0, len(embs))

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
		tsnes = append(tsnes, v1.Embedding{
			Values:   d.RawRowView(i),
			Metadata: maps.Clone(embs[i].Metadata),
		})
	}

	return tsnes, nil
}
