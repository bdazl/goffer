package scenes

import (
	"image"
	"math/rand"

	"github.com/HexHacks/goffer/pkg/animation/epicycles"
)

type EpiOpti struct {
	epi *epicycles.Epicycles
}

func (e *EpiOpti) Init() {
	const (
		ptCnt = 200
		cCnt  = 21
	)

	c := make([]complex128, cCnt)
	for i := range c {
		a, b := rand.Float64(), rand.Float64()
		c[i] = complex(a*0.3-0.15, b*0.3-0.15)
	}

	e.epi = epicycles.New(c, ptCnt)
}

func (e *EpiOpti) Frame(t float64) image.Image {
	return e.epi.Frame(t)
}
