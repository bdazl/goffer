package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	TwoPi = math.Pi * 2.0
)

var (
	Output     = "out/out.gif"
	FPS        = 24
	FrameCount = FPS * 5
	Width      = 512
	Height     = 512
	Total      = 5.

	// https://colorhunt.co/palette/177866
	Palette = color.Palette{
		color.RGBA{0x20, 0x40, 0x51, 0xff},
		color.RGBA{0x3B, 0x69, 0x78, 0xff},
		color.RGBA{0x84, 0xA9, 0xAC, 0xff},
		color.RGBA{0xCA, 0xE8, 0xD5, 0xff},
	}
)

func main() {
	flag.StringVar(&Output, "output", Output, "output file")
	flag.IntVar(&FPS, "fps", FPS, "frames per second")
	flag.IntVar(&FrameCount, "fcount", FrameCount, "frame count")
	flag.Parse()

	Total = float64(FrameCount) / float64(FPS)

	// delay is per frame, in 100ths of a second
	imgs := animate(FrameCount, FPS)

	jiffy := &gif.GIF{
		Image: imgs,
		Delay: getDelays(len(imgs)),
	}

	ofile, err := os.Create(Output)
	if err != nil {
		panic(err)
	}

	err = gif.EncodeAll(ofile, jiffy)
	if err != nil {
		panic(err)
	}
}

func frame(t float64) *image.Paletted {
	var (
		w, h   = float64(Width), float64(Height)
		cx, cy = w / 2.0, h / 2.0
		bounds = image.Rect(0, 0, Width, Height)
	)

	// Initialize the graphic context on an RGBA image
	full := image.NewRGBA(bounds)
	gc := draw2dimg.NewGraphicContext(full)

	// Set some properties
	gc.SetFillColor(Palette[1])
	gc.SetStrokeColor(Palette[2])
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

	return toPaletted(full, Palette)
}

func animate(count int, fps int) []*image.Paletted {
	var (
		ffps = float64(fps)
		out  = make([]*image.Paletted, count)
	)

	for i := 0; i < count; i++ {
		t := float64(i) / ffps
		fmt.Printf("%.3fs\n", t)
		out[i] = frame(t)
	}

	return out
}

func toPaletted(img image.Image, palette color.Palette) *image.Paletted {
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

func getDelays(count int) []int {
	delay := 1
	out := make([]int, count)
	for i := range out {
		out[i] = delay
	}
	return out
}
