package main

import (
	"image"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"
)

func frameFulkonstEtt(t float64) *image.Paletted {
	var (
		palette = Palette1
		w, h    = float64(Width), float64(Height)
		cx, cy  = w / 2.0, h / 2.0
		bounds  = image.Rect(0, 0, Width, Height)
	)

	// Initialize the graphic context on an RGBA image
	img := image.NewRGBA(bounds)
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetFillColor(palette[1])
	gc.SetStrokeColor(palette[2])
	gc.SetLineWidth(2)

	rad := TwoPi * t / Total
	amp := w / 2.0
	cosp, cosc := amp*math.Cos(rad)+cx, amp*math.Cos(rad/2.0)+cx
	sinp, sinc := amp*math.Sin(rad)+cy, amp*math.Sin(rad/2.0)+cy

	// Draw a closed shape
	gc.BeginPath()
	gc.MoveTo(cx, cy)
	gc.LineTo(w, h/2.0)
	gc.QuadCurveTo(cosc, sinc, cosp, sinp)
	gc.Close()
	gc.FillStroke()

	return gifEncodeFrame(img, palette)
}
