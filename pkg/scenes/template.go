package scenes

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/global"
	jimage "github.com/HexHacks/goffer/pkg/image"
	"github.com/HexHacks/goffer/pkg/image/mask"
)

type Template struct {
}

func (_ *Template) Init() {
}

func (_ *Template) Frame(t float64) image.Image {
	var (
		W  = global.Width
		H  = global.Height
		CX = W / 2
		CY = H / 2
		//(C  = image.Point{CX, CY}

		blue        = color.RGBA{0, 0, 255, 255}
		red         = color.RGBA{220, 10, 10, 255}
		uniformBlue = &image.Uniform{blue}
	)

	img, _ := jimage.New()

	// Background: Transparent
	draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)

	srcR := image.Rect(CX, CY, CX+40, CY+50)
	dstP := image.Point{CX / 2, CY / 2}
	dstR := srcR.Sub(srcR.Min).Add(dstP)

	// Draw a red rectangle
	draw.Draw(img, dstR, &image.Uniform{red}, srcR.Min, draw.Src)

	r := 40
	p := image.Point{CX, CY}
	circMask := &mask.Circle{P: p, R: r}

	// Draw blue circle
	draw.DrawMask(
		img, img.Bounds(), // To output image
		uniformBlue, image.ZP, // From blue color
		circMask, image.ZP, // Mask covers all image
		draw.Over)

	return img
}

func (_ *Template) Update(t float64) {
}
