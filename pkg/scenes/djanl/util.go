package djanl

import (
	"image"
	"math"
	"math/rand"

	"github.com/HexHacks/goffer/pkg/math/float"
)

func rev(a []float64) []float64 {
	l := len(a)
	out := make([]float64, l)
	for i, j := 0, l-1; i < l; i++ {
		out[i] = a[j]
		j--
	}
	return out
}

func revc(a []complex128) []complex128 {
	l := len(a)
	out := make([]complex128, l)
	for i, j := 0, l-1; i < l; i++ {
		out[i] = a[j]
		j--
	}
	return out
}

func flattenPts(args [][]complex128) []complex128 {
	out := make([]complex128, 0)
	for _, lst := range args {
		out = append(out, lst...)
	}
	return out
}

// random intervall
func randI(a, b float64) float64 {
	r := rand.Float64()
	atb := b - a
	return a + (atb * r)
}

func randPoint(max image.Point) image.Point {
	r0, r1 := rand.Int(), rand.Int()
	return image.Point{
		r0 % max.X,
		r1 % max.Y,
	}
}

func randComplexPoint(max image.Point) complex128 {
	r0, r1 := rand.Float64(), rand.Float64()
	return complex(
		r0*float64(max.X),
		r1*float64(max.Y),
	)
}

func cmplxCircle(cnt image.Point, a, r float64) complex128 {
	circle := complex(
		math.Cos(a)*r,
		math.Sin(a)*r,
	)

	return PtoC(cnt) + circle
}

func lissajous(cnt image.Point, x, A, B, a, b, δ float64) complex128 {
	lis := complex(
		A*math.Sin(a*x+δ),
		B*math.Sin(b*x),
	)

	return PtoC(cnt) + lis
}

func beatFunc(t float64) float64 {
	const (
		T = tempoPeriod

		a = 1
		b = 0
		//c = 18
		c = 20

		m = T

		cmin = 18
		cmax = 200
	)
	g := func(x float64) float64 {
		return float.Gaussian(x, a, b, c)
	}

	m0 := math.Mod(t, m)
	m1 := math.Mod(-t, m)
	return float.Clamp(g(m0)+g(m1), 0.0, 1.0)
}

func PT(x, y int) image.Point {
	return image.Point{x, y}
}

func PTs(n int) image.Point {
	return image.Point{n, n}
}

func PtoC(p image.Point) complex128 {
	return complex(float64(p.X), float64(p.Y))
}

func CToP(p complex128) image.Point {
	return image.Point{
		X: int(real(p)),
		Y: int(imag(p)),
	}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
