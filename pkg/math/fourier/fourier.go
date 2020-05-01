package fourier

import (
	"math"
	"math/cmplx"

	jcmplx "github.com/HexHacks/goffer/pkg/math/cmplx"
)

const (
	T    = math.Pi
	T2   = 2.0 * math.Pi
	InvT = 1.0/T2 + 0i
)

// Coefficients of the complex fourier series.
// order determines the amount of coefficients returned
// which is the nearest odd number
func Coefficients(path []complex128, order int) []complex128 {
	k := 2*order + 1 // make sure we have an odd number of coefficients
	kh := k / 2

	// should be T / (len(path)-1), however this is also canceled out
	// in the function by the devision of the period
	dt := 1.0 / float64(len(path)-1)

	// used to calculate function for integral
	f := make([]complex128, len(path))

	coef := make([]complex128, k)

	for n := range coef {
		// calculate function for this coef
		for fi := range f {
			t := 2*T*(dt*float64(fi)) - T
			nf := float64(n - kh)
			jnot := complex(0.0, -t*nf)
			f[fi] = path[fi] * cmplx.Exp(jnot)
		}

		coef[n] = InvT * jcmplx.Integrate(-T, T, f)
	}

	return coef
}

// P goes from the frequency domain (coefficients) to the time domain.
// t should be in [0, 1]
func P(t float64, coef []complex128) complex128 {
	// set domain to [-T, T]
	t = T2*t - T

	h := len(coef) / 2
	out := 0 + 0i
	for i, c := range coef {
		n := float64(i - h)
		out += c * cmplx.Exp(complex(0, n*t))
	}
	return out
}
