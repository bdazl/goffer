package main

import (
	"image/color"
)

var (
	Palette = expandPalette(Palette1, 2)

	// https://colorhunt.co/palette/177866
	Palette1 = color.Palette{
		color.RGBA{0x20, 0x40, 0x51, 0xff},
		color.RGBA{0x3B, 0x69, 0x78, 0xff},
		color.RGBA{0x84, 0xA9, 0xAC, 0xff},
		color.RGBA{0xCA, 0xE8, 0xD5, 0xff},
	}
)

func expandPalette(p color.Palette, cnt int) color.Palette {
	np := make(color.Palette, 0, len(p)*cnt*2*2*2)
	for _, c := range p {
		np = append(np, c)
		for i := 0; i < cnt; i++ {
			var m uint8 = uint8(i) * 5
			np = append(np,
				pNextCol(c, m, 0, 0),
				pNextCol(c, m, m, 0),
				pNextCol(c, m, 0, m),
				pNextCol(c, 0, m, 0),
				pNextCol(c, 0, m, m),
				pNextCol(c, 0, 0, m),
				pNextCol(c, m, m, m),

				pNextCol(c, -m, 0, 0),
				pNextCol(c, -m, -m, 0),
				pNextCol(c, -m, 0, -m),
				pNextCol(c, 0, -m, 0),
				pNextCol(c, 0, -m, -m),
				pNextCol(c, 0, 0, -m),
				pNextCol(c, -m, -m, -m),
			)
		}
	}
	return np
}

func pNextCol(c color.Color, r, g, b uint8) color.Color {
	rgba := c.(color.RGBA)
	rgba.R = rgba.R + r
	rgba.G = rgba.G + g
	rgba.B = rgba.B + b
	return rgba
}
