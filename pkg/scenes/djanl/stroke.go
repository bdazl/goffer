package djanl

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/bezier"
)

type stroke struct {
	brush
	curve bezier.Curve
}

func newStroke(img image.Image, pts []complex128) stroke {
	return stroke{
		brush: newBrush(img),
		curve: bezier.New(pts...),
	}
}

func (s *stroke) Draw(dst draw.Image, t float64) {
	ptc := s.curve.Point(t)
	pt := image.Point{int(real(ptc)), int(imag(ptc))}
	s.brush.Draw(dst, pt)
}
