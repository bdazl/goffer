package main

import (
	"math"

	"gonum.org/v1/gonum/spatial/r2"
)

const (
	TwoPi = math.Pi * 2.0
)

func V(x, y float64) r2.Vec {
	return r2.Vec{X: x, Y: y}
}

// [0, 1], radius, translation -> circle
func jexp(f, r, cx, cy float64) (float64, float64) {
	return r*math.Cos(TwoPi*f) + cx, -r*math.Sin(TwoPi*f) + cy
}

func jexpV(f, r, cx, cy float64) r2.Vec {
	x, y := jexp(f, r, cx, cy)
	return r2.Vec{X: x, Y: y}
}

// MorletReal returns the real part of the morlet wavelet
// Sigma is nice
func MorletReal(sigma, t float64) float64 {
	const (
		m12 = -1 / 2
		m14 = -1 / 4
		m34 = -3 / 4
	)
	var (
		s     = sigma
		ss    = s * s
		tt    = t * t
		pim14 = math.Pow(math.Pi, m14)
	)

	c := math.Pow(1.0+exp(-ss)-2*exp(m34*ss), m12)
	k := exp(m12 * ss)
	return c * pim14 * exp(m12*tt) * (math.Cos(s*t) - k)
}

// MorletBnd gives a nicely bounded morlet
// (sigma should be around 2 - 10, but could be anything)
func MorletBnd(sigma, t float64) float64 {
	const (
		a = 1.3 // gives the amplitude bounds of just under 1
		b = 3.5 // the full "range" is [-1, 1]
	)
	return a * MorletReal(sigma, t*b)
}

func exp(x float64) float64 {
	return math.Exp(x)
}

func smoothstep(edge0, edge1, x float64) float64 {
	// Scale, and clamp x to 0..1 range
	x = clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	// Evaluate polynomial
	return x * x * x * (x*(x*6-15) + 10)
}

func clamp(x, low, hi float64) float64 {
	if x < low {
		return low
	}
	if x > hi {
		return hi
	}
	return x
}

// VECTORS

func norm(v r2.Vec) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func length(a, b r2.Vec) float64 {
	return norm(b.Sub(a))
}

func normalize(a r2.Vec) r2.Vec {
	return a.Scale(1.0 / norm(a))
}

// Complex functions
// w: width of image
// h: height of image
// zw: width of complex region to stretch
// zh: height of complex region to stretch
func complexToImage(c complex128, w, h, zw, zh float64) (float64, float64) {
	return (real(c) + zw) * w / (2.0 * zw), (-imag(c) + zh) * h / (2.0 * zh)
}
