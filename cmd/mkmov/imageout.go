package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"time"
)

func pngOutputDir(imgDir string, imgs []image.Image) {
	err := os.MkdirAll(imgDir, 0775)
	panicOn(err)

	times := make([]float64, len(imgs))

	for i, img := range imgs {
		fileName := path.Join(imgDir, fmt.Sprintf("%v.png", i))

		start := time.Now()

		fil, err := os.Create(fileName)
		panicOn(err)

		err = png.Encode(fil, img)
		panicOn(err)

		meas := time.Since(start)
		ms := getMs(meas)
		times[i] = ms

		// fmt.Printf("encoded: %v, time: %.2f\n", fileName, ms)
		panicOn(fil.Close())
	}

	fmt.Printf("png enc stats; ")
	printStats(times)
}
