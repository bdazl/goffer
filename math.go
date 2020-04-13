package main

import (
	"math"

	"gonum.org/v1/gonum/spatial/r2"
)

const (
	TwoPi = math.Pi * 2.0
)

// [0, 1], radius, translation -> circle
func jexp(f, r, cx, cy float64) (float64, float64) {
	return r*math.Cos(TwoPi*f) + cx, -r*math.Sin(TwoPi*f) + cy
}

func jexpV(f, r, cx, cy float64) r2.Vec {
	x, y := jexp(f, r, cx, cy)
	return r2.Vec{X: x, Y: y}
}

func norm(v r2.Vec) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func length(a, b r2.Vec) float64 {
	return norm(b.Sub(a))
}

func normalize(a r2.Vec) r2.Vec {
	return a.Scale(1.0 / norm(a))
}
