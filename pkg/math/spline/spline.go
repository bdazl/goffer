package spline

import (
	"github.com/cnkei/gospline"
)

type Spline struct {
	X, Y gospline.Spline
	Pts  []complex128 // debug
}

func New(pts []complex128) Spline {
	var (
		plen  = len(pts)
		pmaxf = float64(plen - 1)
		X     = make([]float64, plen)
		Y     = make([]float64, plen)
		C     = make([]float64, plen)
	)
	for i, p := range pts {
		t := float64(i) / pmaxf

		X[i] = real(p)
		Y[i] = imag(p)
		C[i] = t
	}

	return Spline{
		X:   gospline.NewCubicSpline(C, X),
		Y:   gospline.NewCubicSpline(C, Y),
		Pts: pts,
	}
}

// Unit interval [0, 1]
func (s *Spline) At(x float64) complex128 {
	return complex(
		s.X.At(x),
		s.Y.At(x),
	)
}

// Range returns interpolated values in [start, end] with step
func (s *Spline) Range(start, end, step float64) []complex128 {
	var (
		XR = s.X.Range(start, end, step)
		YR = s.Y.Range(start, end, step)
	)

	if len(XR) != len(YR) {
		panic("len(XR) != len(YR)")
	}

	out := make([]complex128, len(XR))
	for i, x := range XR {
		y := YR[i]
		out[i] = complex(x, y)
	}

	return out
}
