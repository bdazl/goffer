package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
	"path"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func main() {
	flag.IntVar(&FPS, "fps", FPS, "frames per second")
	flag.IntVar(&FrameCount, "fcount", FrameCount, "frame count")
	flag.Parse()

	rand.Seed(19901231)
	initGlobals()

	imgs := animate(FrameCount, FPS)

	jiffy := &gif.GIF{
		Image: imgs,
		Delay: getDelays(len(imgs), FPS),
	}

	gifOutputFile(path.Join("out", fmt.Sprintf("%v.gif", Pstr)), jiffy)
}

func animate(count int, fps int) []*image.Paletted {
	var (
		ffps = float64(fps)
		out  = make([]*image.Paletted, count)
	)

	t0 := time.Now()
	P.Init()
	fmt.Printf("init time: %.3fms\n", getMs(time.Since(t0)))

	t1 := time.Now()
	times := make([]float64, count)
	for i := 0; i < count; i++ {
		t := float64(i) / ffps

		start := time.Now()
		out[i] = P.Frame(t)
		meas := time.Since(start)

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("seek: %.3fs, build time: %.3fms\n", t, ms)
	}

	fmt.Printf("total build time: %.3fms\n", getMs(time.Since(t1)))
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
	sum := floats.Sum(s)
	avg, std := stat.MeanStdDev(s, nil)
	fmt.Printf("sum: %.3f, μ: %.3f, σ: %.3f (95%% aka ±2σ = ±%.3f)\n", sum, avg, std, 3*std)
}

func gifEncodeFrame(img image.Image, palette color.Palette) *image.Paletted {
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
	ofile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file:", filename)
	err = gif.EncodeAll(ofile, jiffy)
	if err != nil {
		panic(err)
	}
}
