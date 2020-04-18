package float

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("vim-go")
}

// MorletReal returns the real part of the morlet wavelet
// Sigma is nice
func Morlet(sigma, t float64) float64 {
	const (
		m12 = -1.0 / 2.0
		m14 = -1.0 / 4.0
		m34 = -3.0 / 4.0
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
	return a * Morlet(sigma, t*b)
}

func Smoothstep(edge0, edge1, x float64) float64 {
	// Scale, and clamp x to 0..1 range
	x = Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	// Evaluate polynomial
	return x * x * x * (x*(x*6-15) + 10)
}

func Clamp(x, low, hi float64) float64 {
	if x < low {
		return low
	}
	if x > hi {
		return hi
	}
	return x
}

// The parameter a is the height of the curve's peak, b is the position of the center of the peak and c (the standard deviation, sometimes called the Gaussian RMS width)
func Gaussian(x, a, b, c float64) float64 {
	q := x - b
	return a * exp(-q*q/2*c*c)
}

// helper
func exp(x float64) float64 {
	return math.Exp(x)
}
