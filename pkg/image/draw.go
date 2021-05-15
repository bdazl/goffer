package image

import (
	"image"

	"github.com/HexHacks/goffer/pkg/coordsys"
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

// Closed as in, if count is len(pts) then also close the curve
func DrawClosedLines(gc *draw2dimg.GraphicContext, pts []complex128, count int) {
	if count > len(pts) {
		count = len(pts)
	}

	// draw line through all points
	start := coordsys.UnitToImgC(pts[0])
	p := start
	gc.MoveTo(real(p), imag(p))
	for i := 0; i < count; i++ {
		p := coordsys.UnitToImgC(pts[i])
		gc.LineTo(real(p), imag(p))
	}

	if count == len(pts) {
		gc.LineTo(real(start), imag(start))
	}

	gc.Stroke()
}

func DrawLines(gc *draw2dimg.GraphicContext, pts []complex128, count int) {
	if len(pts) < 1 {
		return
	}

	if count > len(pts) {
		count = len(pts)
	}

	// draw line through all points
	start := coordsys.UnitToImgC(pts[0])
	p := start
	gc.MoveTo(real(p), imag(p))
	for i := 0; i < count; i++ {
		p := coordsys.UnitToImgC(pts[i])
		gc.LineTo(real(p), imag(p))
	}

	gc.Stroke()
}

func DrawLinesImgCoords(gc *draw2dimg.GraphicContext, pts []complex128, count int) {
	if len(pts) < 1 {
		return
	}

	if count > len(pts) {
		count = len(pts)
	}

	// draw line through all points
	start := pts[0]
	p := start
	gc.MoveTo(real(p), imag(p))
	for i := 0; i < count; i++ {
		p := pts[i]
		gc.LineTo(real(p), imag(p))
	}

	gc.Stroke()
}
