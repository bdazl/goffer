package djanl

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/bezier"
	"github.com/HexHacks/goffer/pkg/math/spline"
)

type stroke struct {
	brush
	spline spline.Spline
	curve  bezier.Curve // obsolete?
}

func newStroke(img image.Image, pts []complex128) stroke {
	return stroke{
		brush:  newBrush(img),
		spline: spline.New(pts),
		curve:  bezier.New(pts...),
	}
}

func (s *stroke) Range(start, end, step float64) []complex128 {
	return s.spline.Range(start, end, step)
}

func (s *stroke) DrawAt(dst draw.Image, pt complex128) {
	ipt := image.Point{int(real(pt)), int(imag(pt))}
	s.brush.Draw(dst, ipt)
}

// Obsolete ?
func (s *stroke) Draw(dst draw.Image, t float64) {
	ptc := s.curve.Point(t)
	s.DrawAt(dst, ptc)
}
