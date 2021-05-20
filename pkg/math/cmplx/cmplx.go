package cmplx

import (
	"math/rand"

	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/spatial/r2"
)

// Complex functions
// w: width of image
// h: height of image
// zw: width of complex region to stretch
// zh: height of complex region to stretch
func ToImage(c complex128, w, h, zw, zh float64) (float64, float64) {
	return (real(c) + zw) * w / (2.0 * zw), (-imag(c) + zh) * h / (2.0 * zh)
}

func ToVec(c complex128) r2.Vec {
	return r2.Vec{
		X: real(c),
		Y: imag(c),
	}
}

// Integrate a complex valued function with a real domain
func Integrate(a, b float64, c []complex128) complex128 {
	if b < a {
		panic("bad input")
	}

	dx := (b - a) / float64(len(c)-1)

	x := make([]float64, len(c))
	u := make([]float64, len(c))
	v := make([]float64, len(c))

	for i := range x {
		x[i] = a + dx*float64(i)
		u[i] = real(c[i])
		v[i] = imag(c[i])
	}

	return complex(integrate.Trapezoidal(x, u), integrate.Trapezoidal(x, v))
}

func RandomSlice(count int, h float64) []complex128 {
	hh := h / 2.0
	out := make([]complex128, count)
	for i := range out {
		a, b := rand.Float64(), rand.Float64()
		out[i] = complex(a*h-hh, b*h-hh)
	}
	return out
}

func Scale(dst []complex128, s complex128) {
	for i := range dst {
		dst[i] *= s
	}
}

// dst = dst + src
func Add(dst, src []complex128) {
	if len(dst) != len(src) {
		panic("bad lengths")
	}

	for i := range dst {
		dst[i] += src[i]
	}
}

// dst = dst - src
func Sub(dst, src []complex128) {
	if len(dst) != len(src) {
		panic("bad lengths")
	}

	for i := range dst {
		dst[i] -= src[i]
	}
}

func CenterOfMass(pts []complex128) complex128 {
	accum := complex(0, 0)
	for _, p := range pts {
		accum += p
	}

	plen := float64(len(pts))
	return complex(
		real(accum)/plen,
		imag(accum)/plen,
	)
}
