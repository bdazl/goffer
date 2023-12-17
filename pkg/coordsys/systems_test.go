package coordsys

import (
	"testing"

	"github.com/bdazl/goffer/pkg/global"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/spatial/r2"
)

const (
	eps = 2e-10
)

func TestImgAndCdom(t *testing.T) {
	// unit coords
	u := []r2.Vec{
		r2.Vec{X: -1.0, Y: 1.0},
		r2.Vec{X: 1.0, Y: 1.0},
		r2.Vec{X: 1.0, Y: -1.0},
		r2.Vec{X: -1.0, Y: -1.0},
	}
	i := []r2.Vec{
		r2.Vec{X: 0.0, Y: 0.0},
		r2.Vec{X: global.W, Y: 0.0},
		r2.Vec{X: global.W, Y: global.H},
		r2.Vec{X: 0.0, Y: global.H},
	}

	// check boundaries
	for j, unit := range u {
		inEpsilon(t, unit, ImgToUnit(i[j]))
	}

	// ImgToUnit and UnitToImg should be inverses of each other
	for _, img := range i {
		inEpsilon(t, img, UnitToImg(ImgToUnit(img)))
	}
}

func inEpsilon(t *testing.T, exp, res r2.Vec) {
	assert.InDelta(t, exp.X, res.X, eps)
	assert.InDelta(t, exp.Y, res.Y, eps)
}
