package djanl

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/bdazl/goffer/pkg/bezier"
	"github.com/bdazl/goffer/pkg/math/spline"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/exp/rand"
)

type stroke struct {
	brush
	spline spline.Spline
	curve  bezier.Curve // obsolete?
	light0 colorful.Color
	dark0  colorful.Color
}

func newStroke(img image.Image, pts []complex128) stroke {
	// Colors
	lgt := colorful.HappyColor()
	h, s, v := lgt.Hsv()

	rh, rs := rand.Float64()*360.0, rand.Float64()*0.5-0.25

	drkr := colorful.Hsv(h+rh, s*rs, v*0.5)

	return stroke{
		brush:  newBrush(img),
		spline: spline.New(pts),
		curve:  bezier.New(pts...),
		light0: lgt,
		dark0:  drkr,
	}
}

func (s *stroke) Range(start, end, step float64) []complex128 {
	return s.spline.Range(start, end, step)
}

func (s *stroke) DrawColAt(dst draw.Image, pt complex128, col color.Color) {
	ipt := CToP(pt)
	s.brush.DrawColor(dst, ipt, col)
}

func (s *stroke) DrawAt(dst draw.Image, pt complex128) {
	ipt := CToP(pt)
	s.brush.Draw(dst, ipt)
}

func (s *stroke) Draw(dst draw.Image, t float64) {
	ptc := s.curve.Point(t)
	s.DrawAt(dst, ptc)
}
