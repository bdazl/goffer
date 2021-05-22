package djanl

// Made 2021-05-08

import (
	"image"
	"image/color"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/HexHacks/goffer/pkg/bezier"
	"github.com/HexHacks/goffer/pkg/global"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Real value
	// bpm = 100.0
	bpm         = 100.0
	tempoFreq   = bpm / 60.0
	tempoPeriod = 60.0 / bpm

	cutoutCnt    = 20
	bezierPoints = 600

	twoPi    = math.Pi * 2.0
	piHalf   = math.Pi / 2.0
	piFourth = math.Pi / 4.0

	cutout = 100 // px i bild
	cutC   = 15.0

	// because I was stupid and exported large images for a long time
	// I need a conversion factor on some values
	correct = 1.0 / 2048.0
)

var (
	Width  = global.Width
	Height = global.Height
	W      = float64(Width)
	H      = float64(Height)

	CX = Width / 2
	CY = Height / 2

	//(C  = image.Point{CX, CY}
	Dur = global.Total

	MaxTime = float64(global.FrameCount-1) / float64(global.FPS)

	cutoutR = W / cutC

	blue        = color.RGBA{0, 0, 255, 255}
	red         = color.RGBA{220, 10, 10, 255}
	uniformBlue = &image.Uniform{blue}

	//bezierPoints = int(math.Ceil(Dur)) * 2
)

func resetGlobals() {
	Width = global.Width
	Height = global.Height

	W = float64(Width)
	H = float64(Height)

	CX = Width / 2
	CY = Height / 2

	cutoutR = W / cutC

	Dur = global.Total
	MaxTime = float64(global.FrameCount-1) / float64(global.FPS)

	if MaxTime <= 0.0 {
		MaxTime = 1
	}
}

type Djanl struct {
	palette []colorful.Color
	refImgs []refImage
	strokes []stroke
	bgCurve bezier.Curve
}

type ImageSub interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}
