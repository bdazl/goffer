package djanl

// Made 2021-05-08

import (
	"image"
	"image/color"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/HexHacks/goffer/pkg/global"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Real value
	// bpm = 100.0
	bpm         = 100.0
	tempoFreq   = bpm / 60.0
	tempoPeriod = 60.0 / bpm

	cutoutCnt = 20
)

var (
	W  = global.Width
	H  = global.Height
	CX = W / 2
	CY = H / 2
	//(C  = image.Point{CX, CY}
	Dur = global.Total

	blue        = color.RGBA{0, 0, 255, 255}
	red         = color.RGBA{220, 10, 10, 255}
	uniformBlue = &image.Uniform{blue}

	cutoutR = image.Rect(0, 0, 100, 100)

	//bezierPoints = int(math.Ceil(Dur)) * 2
	bezierPoints = 500
	twoPi        = math.Pi * 2.0
)

type Djanl struct {
	palette []colorful.Color
	refImgs []refImage
	strokes []stroke
}

type ImageSub interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}
