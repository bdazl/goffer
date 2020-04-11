package main

import (
	"flag"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
)

var (
	output = "out/hello.png"

	// https://colorhunt.co/palette/177866
	Palette = color.Palette{
		color.RGBA{0x20, 0x40, 0x51, 0xff},
		color.RGBA{0x3B, 0x69, 0x78, 0xff},
		color.RGBA{0x84, 0xA9, 0xAC, 0xff},
		color.RGBA{0xCA, 0xE8, 0xD5, 0xff},
	}
)

func main() {
	flag.StringVar(&output, "output", output, "output file")
	flag.Parse()

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 512, 512))
	gc := draw2dimg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(Palette[1])
	gc.SetStrokeColor(Palette[2])
	gc.SetLineWidth(2)

	// Draw a closed shape
	gc.BeginPath()    // Initialize a new path
	gc.MoveTo(10, 10) // Move to a position to start the new path
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()

	// Save to file
	draw2dimg.SaveToPngFile(output, toPaletted(dest, Palette))
}

func toPaletted(img image.Image, palette color.Palette) image.Image {
	bnds := img.Bounds()
	out := image.NewPaletted(bnds, palette)

	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			idx := palette.Index(img.At(x, y))
			out.SetColorIndex(x, y, uint8(idx))
		}
	}
	return out
}
