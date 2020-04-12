package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"time"

	"github.com/gonum/stat"
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

	imgs := animate(FrameCount, FPS)

	jiffy := &gif.GIF{
		Image: imgs,
		Delay: getDelays(len(imgs), FPS),
	}

	gifOutputFile(Output, jiffy)
}

func animate(count int, fps int) []*image.Paletted {
	var (
		ffps = float64(fps)
		out  = make([]*image.Paletted, count)
	)

	times := make([]float64, count)
	for i := 0; i < count; i++ {
		t := float64(i) / ffps

		start := time.Now()
		out[i] = gifFrameEncode(Frame(t), Palette)
		meas := time.Since(start)

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("seek: %.3fs, build time: %.3fms\n", t, ms)
	}

	printStats(times)

	return out
}

func getMs(dur time.Duration) float64 {
	return float64(dur.Nanoseconds()) * 1e-6
}

func getDelays(count, fps int) []int {
	// delay is per frame, in 100ths of a second
	delay := 100 / fps

	out := make([]int, count)
	for i := range out {
		out[i] = delay
	}
	return out
}

func printStats(s []float64) {
	avg, std := stat.MeanStdDev(s, nil)
	fmt.Printf("μ: %.3f, σ: %.3f (95%% aka ±2σ = ±%.3f)\n", avg, std, 3*std)
}

func gifFrameEncode(img image.Image, palette color.Palette) *image.Paletted {
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

func gifOutputFile(filename string, jiffy *gif.GIF) {
	ofile, err := os.Create(Output)
	if err != nil {
		panic(err)
	}

	err = gif.EncodeAll(ofile, jiffy)
	if err != nil {
		panic(err)
	}
}
