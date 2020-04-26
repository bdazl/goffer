package svg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/spatial/r2"
)

const (
	eps = 2e-10
)

func TestParsePath(t *testing.T) {
	const (
		simplePath = `
m 10.0,20.0 c 0,0 10.0,20.0 -10.0,20.0 -10.0,-20.0
`
	)

	_, err := parsePath("")
	assert.Error(t, err)

	exp := []Operation{
		Operation{
			Type:   Move,
			Points: []Point{Point{10.0, 20.0}},
		},
		Operation{
			Type: Curve,
			Points: []Point{
				Point{0.0, 0.0},
				Point{10.0, 20.0},
				Point{-10.0, 20.0},
				Point{-10.0, -20.0},
			},
		},
	}

	ops, err := parsePath(simplePath)
	assert.NoError(t, err)
	assertOps(t, exp, ops)

	upperPath := strings.ToUpper(simplePath)

	ops, err = parsePath(upperPath)
	assert.NoError(t, err)
	assertOps(t, exp, ops)
}

func assertOps(t *testing.T, exp, out []Operation) {
	assert.Equal(t, len(exp), len(out))

	for i, ee := range exp {
		oo := out[i]

		assert.Equal(t, ee.Type, oo.Type)
		assert.Equal(t, len(ee.Points), len(oo.Points))

		for pi, ep := range ee.Points {
			op := oo.Points[pi]
			inEpsilon(t, toVec(ep), toVec(op))
		}
	}
}

func toVec(p Point) r2.Vec {
	return r2.Vec{X: p.X, Y: p.Y}
}

func inEpsilon(t *testing.T, exp, res r2.Vec) {
	assert.InDelta(t, exp.X, res.X, eps)
	assert.InDelta(t, exp.Y, res.Y, eps)
}
