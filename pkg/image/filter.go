package image

import (
	"image"
	"image/color"
)

type FilterFunc = func(x, y int, col color.Color) color.Color

type Filter struct {
	Img image.Image
	FilterFunc
}

func (c *Filter) ColorModel() color.Model {
	return color.RGBAModel
}

func (f *Filter) Bounds() image.Rectangle {
	return f.Img.Bounds()
}

func (f *Filter) At(x, y int) color.Color {
	col := f.Img.At(x, y)
	return f.FilterFunc(x, y, col)
}
