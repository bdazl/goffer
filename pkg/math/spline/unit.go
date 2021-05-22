package spline

import "github.com/cnkei/gospline"

type Unit struct {
	gospline.Spline
}

func NewUnit(pts []float64) Unit {
	var (
		plen  = len(pts)
		pmaxf = float64(plen - 1)
		C     = make([]float64, plen)
	)

	for i := range pts {
		C[i] = float64(i) / pmaxf
	}

	return Unit{
		Spline: gospline.NewCubicSpline(C, pts),
	}
}

func (u *Unit) At(x float64) float64 {
	return u.Spline.At(x)
}

func (u *Unit) Range(start, end, step float64) []float64 {
	return u.Spline.Range(start, end, step)
}
