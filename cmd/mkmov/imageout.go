package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"time"
)

func pngOutputMany(imgs []image.Image) {
	imgDir := imageDirectory()
	err := os.MkdirAll(imgDir, 0775)
	panicOn(err)

	times := make([]float64, len(imgs))

	for i, img := range imgs {
		start := time.Now()

		outputPng(i, img)

		meas := time.Since(start)
		ms := getMs(meas)
		times[i] = ms

		// fmt.Printf("encoded: %v, time: %.2f\n", fileName, ms)
	}

	fmt.Printf("png enc stats; ")
	printStats(times)
}

func outputPng(idx int, img image.Image) {
	imgDir := imageDirectory()
	fileName := path.Join(imgDir, fmt.Sprintf("%v.png", idx))

	fil, err := os.Create(fileName)
	panicOn(err)
	defer fil.Close()

	err = png.Encode(fil, img)
	panicOn(err)
}
