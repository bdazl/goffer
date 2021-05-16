package djanl

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/bezier"
)

type stroke struct {
	brush
	curve  bezier.Curve
	SX, SY float64
}

func newStroke(img image.Image, pts []complex128, max complex128) stroke {
	return stroke{
		brush: newBrush(img),
		curve: bezier.New(pts...),
		SX:    real(max),
		SY:    imag(max),
	}
}

func (s *stroke) Draw(dst draw.Image, t float64) {
	ndc := s.curve.Point(t)
	ptc := complex(
		real(ndc)*s.SX,
		imag(ndc)*s.SY,
	)

	pt := image.Point{int(real(ptc)), int(imag(ptc))}
	s.brush.Draw(dst, pt)
}
