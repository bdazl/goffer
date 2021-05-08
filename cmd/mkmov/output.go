package main

import (
	"fmt"
	"image"
	"image/gif"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/HexHacks/goffer/pkg/global"
	"github.com/HexHacks/goffer/pkg/palette"
)

type FileType string

const (
	GIF FileType = "gif"
	MP4 FileType = "mp4"
)

func OutputFile(imgs []image.Image) {
	switch OutputFileType {
	case MP4:
		mp4OutputFile(imgs)
	case GIF:
		gifOutputFile(imgs)
	default:
		fmt.Printf("bad output file type")
		os.Exit(1)
	}
}

func backupOld(filename string) {
	if !fileExists(filename) {
		return
	}

	old := filename

	i := 0
	dir := path.Dir(filename)
	fil := path.Base(filename)
	filext := strings.Split(fil, ".")

	for fileExists(filename) {
		id := strconv.Itoa(i)
		filename = path.Join(dir, filext[0]+"_"+id+"."+filext[1])

		i++
	}

	fmt.Printf("backing up file: %v -> %v\n", old, filename)

	src, err := os.Open(old)
	panicOn(err)
	defer src.Close()

	dst, err := os.Create(filename)
	panicOn(err)
	defer dst.Close()

	_, err = io.Copy(dst, src)
	panicOn(err)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func gifOutputFile(imgs []image.Image) {
	filename := videoOutFilename()
	if Backup {
		backupOld(filename)
	}

	fmt.Println("encoding images with gif palette...")
	t0 := time.Now()
	encoded := make([]*image.Paletted, len(imgs))
	for i, im := range imgs {
		encoded[i] = gifEncodeFrame(im)
	}
	fmt.Println("encode time:", time.Since(t0))

	jiffy := &gif.GIF{
		Image: encoded,
		Delay: gifGetDelays(len(imgs), global.FPS),
	}

	ofile, err := os.Create(filename)
	panicOn(err)

	defer ofile.Close()

	fmt.Println("writing to file:", filename)
	err = gif.EncodeAll(ofile, jiffy)
	panicOn(err)
}

func gifEncodeFrame(img image.Image) *image.Paletted {
	bnds := img.Bounds()
	out := image.NewPaletted(bnds, palette.Palette)

	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			idx := palette.Palette.Index(img.At(x, y))
			out.SetColorIndex(x, y, uint8(idx))
		}
	}
	return out
}

func gifGetDelays(count, fps int) []int {
	// delay is per frame, in 100ths of a second
	delay := 100 / fps

	out := make([]int, count)
	for i := range out {
		out[i] = delay
	}
	return out
}
