package r2

import (
	"gonum.org/v1/gonum/spatial/r2"
)

// t [0, 1]
func QuadBezier(t float64, a, b, c r2.Vec) r2.Vec {
	omT := 1.0 - t
	p0 := a.Scale(omT * omT)
	p1 := b.Scale(2 * t * omT)
	p2 := b.Scale(t * t)
	return p0.Add(p1).Add(p2)
}

func CubeBezier(t float64, a, b, c, d r2.Vec) r2.Vec {
	omT := 1.0 - t
	p0 := QuadBezier(t, a, b, c).Scale(omT)
	p1 := QuadBezier(t, b, c, d).Scale(t)
	return p0.Add(p1)
}
