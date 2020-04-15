package main

import (
	"flag"
	"fmt"
	"image"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// cmd line arguments
var (
	FPS        = 30
	FrameCount = FPS * 4
	Width      = 512
	Height     = 512
	Backup     = false
)

func main() {
	flag.IntVar(&FPS, "fps", FPS, "frames per second")
	flag.IntVar(&FrameCount, "fcount", FrameCount, "frame count")
	flag.BoolVar(&Backup, "backup", Backup, "if file exists, do backup")
	flag.StringVar(&ActiveProject, "proj", ActiveProject, "active project")
	flag.Parse()

	rand.Seed(19901231)
	initGlobals()

	imgs := animate(FrameCount, FPS)

	OutputFile(imgs)
}

func animate(count int, fps int) []image.Image {
	t0 := time.Now()
	out := make([]image.Image, count)
	P.Init()
	fmt.Printf("init time: %.3fms\n", getMs(time.Since(t0)))

	t1 := time.Now()
	times := make([]float64, count)

	ffps := float64(fps)
	for i := 0; i < count; i++ {
		t := float64(i) / ffps

		start := time.Now()
		out[i] = P.Frame(t)
		meas := time.Since(start)

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("seek: %.3fs, build time: %.3fms\n", t, ms)
	}

	fmt.Printf("total build time: %v\n", time.Since(t1))
	printStats(times)

	return out
}

func getMs(dur time.Duration) float64 {
	return float64(dur.Nanoseconds()) * 1e-6
}

func printStats(s []float64) {
	sum := floats.Sum(s)
	avg, std := stat.MeanStdDev(s, nil)
	fmt.Printf("sum: %.3f, μ: %.3f, σ: %.3f (95%% aka ±2σ = ±%.3f)\n", sum, avg, std, 3*std)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
