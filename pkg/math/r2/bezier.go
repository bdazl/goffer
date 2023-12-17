package r2

import (
	"gonum.org/v1/gonum/spatial/r2"
)

// t [0, 1]
func QuadBezier(t float64, a, b, c r2.Vec) r2.Vec {
	var (
		omT = 1.0 - t
		p0  = r2.Scale(omT*omT, a)
		p1  = r2.Scale(2*t*omT, b)
		p2  = r2.Scale(t*t, b)
		a0  = r2.Add(p0, p1)
		a1  = r2.Add(a0, p2)
	)
	return a1
}

func CubeBezier(t float64, a, b, c, d r2.Vec) r2.Vec {
	var (
		omT = 1.0 - t
		b0  = QuadBezier(t, a, b, c)
		b1  = QuadBezier(t, b, c, d)
		p0  = r2.Scale(omT, b0)
		p1  = r2.Scale(t, b1)
	)
	return r2.Add(p0, p1)
}
