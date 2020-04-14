package main

import (
	"fmt"
	"image"
	"image/gif"
	"os"
	"time"
)

type FileType int

const (
	GIF = iota + 1
)

func OutputFile(filename string, imgs []image.Image, ft FileType) {
	switch ft {
	case GIF:
		gifOutputFile(filename, imgs)
	default:
		fmt.Printf("bad output file type")
		os.Exit(1)
	}
}

func gifOutputFile(filename string, imgs []image.Image) {
	fmt.Println("encoding images with gif palette...")
	t0 := time.Now()
	encoded := make([]*image.Paletted, len(imgs))
	for i, im := range imgs {
		encoded[i] = gifEncodeFrame(im)
	}
	fmt.Println("encode time:", time.Since(t0))

	jiffy := &gif.GIF{
		Image: encoded,
		Delay: getDelays(len(imgs), FPS),
	}

	ofile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	t1 := time.Now()
	defer fmt.Println("write time:", time.Since(t1))
	defer ofile.Close()

	fmt.Println("writing to file:", filename)
	err = gif.EncodeAll(ofile, jiffy)
	if err != nil {
		panic(err)
	}
}

func gifEncodeFrame(img image.Image) *image.Paletted {
	bnds := img.Bounds()
	out := image.NewPaletted(bnds, Palette)

	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			idx := Palette.Index(img.At(x, y))
			out.SetColorIndex(x, y, uint8(idx))
		}
	}
	return out
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
