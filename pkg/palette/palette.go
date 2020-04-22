package palette

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
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

	Debug = color.Palette{
		rgba("#000000"), // black
		rgba("#ffffff"), // white
		rgba("#ff0000"), // red extreme
		rgba("#00ff00"), // green extreme
		rgba("#0000ff"), // blue extreme
		rgba("#ffff00"), // yellow extreme
		rgba("#ff00ff"), // purple extreme
		rgba("#00ffff"), // cyan extreme

		// red palette
		rgba("#a70000"),
		rgba("#ff0000"),
		rgba("#ff5252"),
		rgba("#ff7b7b"),
		rgba("#ffbaba"),

		// blue palette
		rgba("#005073"),
		rgba("#107dac"),
		rgba("#189ad3"),
		rgba("#1ebbd7"),
		rgba("#71c7ec"),

		// green
		rgba("#006203"),
		rgba("#0f9200"),
		rgba("#30cb00"),
		rgba("#4ae54a"),
		rgba("#a4fba6"),
	}
)

func rgba(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}

func expandPalette(p color.Palette, cnt int) color.Palette {
	np := make(color.Palette, 0, len(p)*cnt*2*2*2)

	// make sure palette is compatible with old stuff
	for _, c := range p {
		np = append(np, c)
	}

	for _, c := range p {
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
