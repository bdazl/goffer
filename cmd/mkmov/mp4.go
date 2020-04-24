package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/HexHacks/goffer/pkg/global"
)

func mp4OutputFile(filename string, imgs []image.Image) {
	// mp4 is created by outputting images as .png files
	// and letting ffmpeg take care of the rest

	// create a directory where the
	baseD := path.Dir(filename)
	out := path.Base(filename)

	dir := path.Join(baseD, ActiveProject+".d")
	err := os.MkdirAll(dir, 0775)
	panicOn(err)

	times := make([]float64, len(imgs))

	for i, img := range imgs {
		fileName := path.Join(dir, fmt.Sprintf("%v.png", i))

		start := time.Now()

		fil, err := os.Create(fileName)
		panicOn(err)
		defer fil.Close()

		err = png.Encode(fil, img)
		panicOn(err)

		meas := time.Since(start)
		ms := getMs(meas)
		times[i] = ms

		fmt.Printf("encoded: %v, time: %.2f\n", fileName, ms)
	}

	fmt.Printf("png enc stats; ")
	printStats(times)

	ffmpeg(dir, out)
}

func ffmpeg(folder, out string) {
	//ffmpeg -r 30 -f image2 -s 512x512 -i $1/%d.png -vcodec libx264 -crf 25  -pix_fmt yuv420p $2

	const (
		crf    = "25" // quality: lower is better, according to docs: "15-25 is usually good"
		codec  = "libx264"
		pixFmt = "yuv420p"
	)
	var (
		frameRate = strconv.Itoa(global.FPS)
		pngFmt    = path.Join(folder, "%d.png")
		dims      = fmt.Sprintf("%vx%v", global.Width, global.Height)
	)

	args := []string{
		"-r", frameRate,
		"-f", "image2",
		"-s", dims,
		"-i", pngFmt,
		"-vcodec", codec,
		"-crf", crf,
		"-pix_fmt", pixFmt,
		path.Join(folder, out),
	}

	start := time.Now()

	fmt.Printf("ffmpeg %v", args)
	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Run()
	panicOn(err)

	meas := time.Since(start)

	fmt.Printf("running ffmpeg took: %v\n", meas)
}
