package image

import (
	"image"

	"github.com/HexHacks/goffer/pkg/global"
	"github.com/HexHacks/goffer/pkg/palette"

	"github.com/llgcode/draw2d/draw2dimg"
)

func New() (*image.RGBA, *draw2dimg.GraphicContext) {
	bounds := image.Rect(0, 0, global.Width, global.Height)
	img := image.NewRGBA(bounds)
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetFillColor(palette.Palette[1])
	gc.SetStrokeColor(palette.Palette[2])
	gc.SetLineWidth(2)

	return img, gc
}
