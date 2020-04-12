package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
)

const (
	TwoPi = math.Pi * 2.0
)

var (
	Frame = frameFulkonstEtt
)

var (
	Output     = "out/out.gif"
	FPS        = 30
	FrameCount = FPS * 4
	Width      = 512
	Height     = 512
)

// calculated
var (
	Total = 5.
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

func animate(count int, fps int) []*image.Paletted {
	var (
		ffps = float64(fps)
		out  = make([]*image.Paletted, count)
	)

	for i := 0; i < count; i++ {
		t := float64(i) / ffps
		fmt.Printf("%.3fs\n", t)
		out[i] = GifEncode(Frame(t), Palette)
	}

	return out
}

func GifEncode(img image.Image, palette color.Palette) *image.Paletted {
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
	delay := 100 / FPS
	out := make([]int, count)
	for i := range out {
		out[i] = delay
	}
	return out
}
