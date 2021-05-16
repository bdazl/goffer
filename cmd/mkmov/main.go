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

type ImageChan = chan image.Image

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

	// image writer thread
	imgCh := make(ImageChan, 3)
	writerDone := make(chan bool, 0)

	go imageWriter(imgCh, writerDone)

	animate(imgCh)

	close(imgCh) // report done to imageWriter
	b, ok := <-writerDone
	if !b || !ok {
		fmt.Errorf("image writer reported error\n")
	}

	outputMov()

	progduration := time.Since(start)
	fmt.Printf("the program runtime was %v\n", progduration)
}

func animate(ch ImageChan) {
	t0 := time.Now()
	P.Init()
	fmt.Printf("init time: %.3fms\n", getMs(time.Since(t0)))

	t1 := time.Now()
	times := make([]float64, global.FrameCount)

	ffps := float64(global.FPS)
	for i := 0; i < global.FrameCount; i++ {
		t := float64(i) / ffps

		start := time.Now()
		img := P.Frame(t)
		meas := time.Since(start)

		// Send image to writer
		ch <- img

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("frame: %v, seek: %.2fs, build time: %.2fms\n", i, t, ms)
	}

	fmt.Printf("total build time: %v\n", time.Since(t1))
	printStats(times)
}

func imageWriter(imgCh ImageChan, done chan bool) {
	i := 0
	for img, ok := <-imgCh; ok; img, ok = <-imgCh {
		outputPng(i, img)
		fmt.Printf("wrote to %v.png\n", i)
		i++
	}

	fmt.Println("image worker done")
	done <- true
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
