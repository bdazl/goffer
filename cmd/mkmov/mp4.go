package main

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/HexHacks/goffer/pkg/global"
)

func mp4OutputFile(imgs []image.Image) {
	// mp4 is created by outputting images as .png files
	// and letting ffmpeg take care of the rest

	// create a directory where the
	out := videoBaseFilename()
	outDir := projOutDirectory()
	imgDir := imageDirectory()

	pngOutputDir(imgDir, imgs)

	outFile := path.Join(outDir, videoBaseFilename())
	if Backup {
		backupOld(outFile)
	} else if fileExists(outFile) {
		_ = os.Remove(outFile)
	}

	ffmpeg(imgDir, outDir, out)
}

func ffmpeg(imgDir, outDir, out string) {
	//ffmpeg -r 30 -f image2 -s 512x512 -i $1/%d.png -vcodec libx264 -crf 25  -pix_fmt yuv420p $2

	const (
		crf    = "20" // quality: lower is better, according to docs: "15-25 is usually good"
		codec  = "libx264"
		pixFmt = "yuv420p"
	)
	var (
		frameRate = strconv.Itoa(global.FPS)
		pngFmt    = path.Join(imgDir, "%d.png")
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
		path.Join(outDir, out),
	}

	start := time.Now()

	fmt.Println("ffmpeg", args)
	cmd := exec.Command("ffmpeg", args...)
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	panicOn(err)

	meas := time.Since(start)

	fmt.Printf("running ffmpeg took: %v\n", meas)
}
