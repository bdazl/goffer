package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
)

func drawCommon(palette color.Palette) (*image.RGBA, *draw2dimg.GraphicContext) {
	bounds := image.Rect(0, 0, Width, Height)
	img := image.NewRGBA(bounds)
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetFillColor(palette[1])
	gc.SetStrokeColor(palette[2])
	gc.SetLineWidth(2)

	return img, gc
}
