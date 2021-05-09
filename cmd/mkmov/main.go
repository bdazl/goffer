package main

import (
	"flag"
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/HexHacks/goffer/pkg/global"
	"github.com/HexHacks/goffer/pkg/scenes"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// cmd line arguments
var (
	Backup        = false
	Verbose       = false
	ActiveProject = scenes.LastScene()
	P             = scenes.Scenes[ActiveProject]
)

var (
	OutputFileType = MP4
)

func main() {
	start := time.Now()

	flag.IntVar(&global.FPS, "fps", global.FPS, "frames per second")
	flag.IntVar(&global.Width, "w", global.Width, "width")
	flag.IntVar(&global.Height, "h", global.Height, "height")
	flag.IntVar(&global.FrameCount, "fcount", global.FrameCount, "frame count")
	flag.BoolVar(&Backup, "backup", Backup, "if file exists, do backup")
	flag.StringVar(&ActiveProject, "proj", ActiveProject, "active project")
	flag.BoolVar(&Verbose, "verbose", Verbose, "more output, notably from ffmpeg")
	flag.Parse()

	rand.Seed(19901231)

	P = scenes.Scenes[ActiveProject]
	global.InitGlobals()

	imgs := animate()

	OutputFile(imgs)

	progduration := time.Since(start)
	fmt.Printf("the program runtime was %v\n", progduration)
}

func animate() []image.Image {
	t0 := time.Now()
	out := make([]image.Image, global.FrameCount)
	P.Init()
	fmt.Printf("init time: %.3fms\n", getMs(time.Since(t0)))

	t1 := time.Now()
	times := make([]float64, global.FrameCount)

	ffps := float64(global.FPS)
	for i := 0; i < global.FrameCount; i++ {
		t := float64(i) / ffps

		start := time.Now()
		out[i] = P.Frame(t)
		meas := time.Since(start)

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("frame: %v, seek: %.3fs, build time: %.3fms\n", i, t, ms)
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
