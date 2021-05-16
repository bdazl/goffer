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
	Parallel      = false
	ActiveProject = scenes.LastScene()
	P             = scenes.Scenes[ActiveProject]
)

var (
	OutputFileType = MP4
)

type Frame struct {
	Number int
	image.Image
}

type ImageChan = chan Frame

func main() {
	start := time.Now()

	flag.IntVar(&global.FPS, "fps", global.FPS, "frames per second")
	flag.IntVar(&global.Width, "w", global.Width, "width")
	flag.IntVar(&global.Height, "h", global.Height, "height")
	flag.IntVar(&global.FrameCount, "fcount", global.FrameCount, "frame count")
	flag.BoolVar(&Backup, "backup", Backup, "if file exists, do backup")
	flag.StringVar(&ActiveProject, "proj", ActiveProject, "active project")
	flag.BoolVar(&Verbose, "verbose", Verbose, "more output, notably from ffmpeg")
	flag.BoolVar(&Parallel, "parallel", Verbose, "frames computed in parallel (THREAD SAFETY)")
	flag.Parse()

	rand.Seed(19901231)

	P = scenes.Scenes[ActiveProject]
	global.InitGlobals()

	renderImages()
	outputMov()

	progduration := time.Since(start)
	fmt.Printf("the program runtime was %v\n", progduration)
}

func renderImages() {
	// image writer thread
	imgCh := make(ImageChan, 30)
	writerDone := make(chan bool, 0)

	go imageWriter(imgCh, writerDone)

	initProject()

	if Parallel {
		renderAllParallel(imgCh)
	} else {
		renderAllSingle(imgCh)
	}

	close(imgCh) // report done to imageWriter
	b, ok := <-writerDone
	if !b || !ok {
		fmt.Printf("image writer reported error\n")
	}
}

func renderAllParallel(ch ImageChan) {
	var (
		workers = 8
		dones   = make([]chan bool, workers)
		fc      = global.FrameCount

		// integer division (some frames might be left out)
		wrkrFrameCount = fc / workers

		// rounding error frames is up to last worker
		slop = fc - wrkrFrameCount*workers
	)

	// bootup workers
	for w := 0; w < workers; w++ {
		dones[w] = make(chan bool, 0)
		start := wrkrFrameCount * w

		count := wrkrFrameCount
		if w == workers-1 {
			count += slop
		}

		fmt.Printf("bootup worker %v\n", w)
		go renderRange(start, count, ch, dones[w])
	}

	// wait for workers
	for i, done := range dones {
		b, ok := <-done
		if !b || !ok {
			pnk := fmt.Sprintf("worker %v failed", i)
			panic(pnk)
		}
	}
}

func renderRange(start, count int, ch ImageChan, done chan bool) {
	ffps := float64(global.FPS)
	for i := start; i < start+count; i++ {
		t := float64(i) / ffps

		start := time.Now()
		img := P.Frame(t)
		meas := time.Since(start)

		// Send image to writer
		ch <- Frame{
			Image:  img,
			Number: i,
		}

		ms := getMs(meas)
		fmt.Printf("frame: %v, seek: %.2fs, build time: %.2fms\n", i, t, ms)
	}

	done <- true
}

func renderAllSingle(ch ImageChan) {

	t1 := time.Now()
	times := make([]float64, global.FrameCount)

	ffps := float64(global.FPS)
	for i := 0; i < global.FrameCount; i++ {
		t := float64(i) / ffps

		start := time.Now()
		img := P.Frame(t)
		meas := time.Since(start)

		// Send image to writer
		ch <- Frame{
			Image:  img,
			Number: i,
		}

		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("frame: %v, seek: %.2fs, build time: %.2fms\n", i, t, ms)
	}

	fmt.Printf("total build time: %v\n", time.Since(t1))
	printStats(times)
}

func imageWriter(imgCh ImageChan, done chan bool) {
	var (
		fc = float64(global.FrameCount)
	)
	defer fmt.Println("image worker done")

	count := 0
	for frame, ok := <-imgCh; ok; frame, ok = <-imgCh {
		outputPng(frame.Number, frame.Image)

		count++
		progress := 100 * float64(count) / fc

		fmt.Printf("%.2f%%: wrote to %v.png\n", progress, frame.Number)
	}

	done <- true
}

func initProject() {
	t0 := time.Now()
	P.Init()
	fmt.Printf("init time: %.3fms\n", getMs(time.Since(t0)))
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
