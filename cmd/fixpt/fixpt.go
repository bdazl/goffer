package main

import (
	"flag"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
)

var (
	output = "out/hello.png"
)

func main() {
	flag.StringVar(&output, "output", output, "output file")
	flag.Parse()

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
	gc := draw2dimg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(2)

	// Draw a closed shape
	gc.BeginPath()    // Initialize a new path
	gc.MoveTo(10, 10) // Move to a position to start the new path
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()

	// Save to file
	draw2dimg.SaveToPngFile(output, dest)
}
