package r2

import (
	"math"

	jmath "github.com/bdazl/goffer/pkg/math"

	"gonum.org/v1/gonum/spatial/r2"
)

func V(x, y float64) r2.Vec {
	return r2.Vec{X: x, Y: y}
}

// [0, 1], radius, translation -> circle
func Exp(f, r, cx, cy float64) (float64, float64) {
	return r*math.Cos(jmath.Tau*f) + cx, -r*math.Sin(jmath.Tau*f) + cy
}

func ExpV(f, r, cx, cy float64) r2.Vec {
	x, y := Exp(f, r, cx, cy)
	return r2.Vec{X: x, Y: y}
}

// VECTORS

func Min(a, b r2.Vec) r2.Vec {
	return r2.Vec{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
	}
}

func Max(a, b r2.Vec) r2.Vec {
	return r2.Vec{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
	}
}

func Norm(v r2.Vec) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func Angle(v r2.Vec) float64 {
	return math.Acos(v.X / Norm(v))
}

func Length(a, b r2.Vec) float64 {
	return Norm(r2.Sub(b, a))
}

func Normalize(a r2.Vec) r2.Vec {
	return r2.Scale(1.0/Norm(a), a)
}

func ToSpherical(rect r2.Vec) r2.Vec {
	return r2.Vec{
		X: Norm(rect),
		Y: Angle(rect),
	}
}
