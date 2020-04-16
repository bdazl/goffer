package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

const (
	importImg = "assets/jagtitt.jpg"
)

type ImgImport struct {
	Img image.Image
}

func (i *ImgImport) Init() {
	fil, err := os.Open(importImg)
	panicOn(err)

	i.Img, err = jpeg.Decode(fil)
	panicOn(err)
}

func (i *ImgImport) getColor(x, y int, t float64) color.Color {
	bnds := i.Img.Bounds()
	//fmt.Printf("start x: %v, y %v, calc: %v\n", x, y, -float64(2*y)/H-1)

	// translate into complex coordinat where
	// [0,0] -> [-1, i]
	// [W,H] -> [1, -1]
	z := complex(float64(2*x)/W-1, -float64(2*y)/H+1)

	// make some transformation
	//w := (z + 2) * (z + 2) * (z - 1 - 2i) * (z + 1i)
	w := 2.0 * z * z * z

	// interpolate between identity and transformation
	o := z + complex(smoothstep(0.0, 1.0, t), 0.0)*(z-w)

	// transform back into image coordinates
	xx, yy := complexToImage(o, W, H, 1.0, 1.0)

	// modulate to make sure we're always inside image
	iw, ih := float64(bnds.Max.X), float64(bnds.Max.Y)
	x = int(math.Mod(xx, iw))
	y = int(math.Mod(yy, ih))

	//fmt.Printf("x: %v, y: %v, z: %v, o: %v, xx: %v, yy: %v\n", x, y, z, o, xx, yy)

	return i.Img.At(x, y)
}

func (i *ImgImport) Frame(t float64) image.Image {
	img, _ := drawCommon(Palette)

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			c := i.getColor(x, y, t)
			img.Set(x, y, c)
		}
	}

	return img
}
