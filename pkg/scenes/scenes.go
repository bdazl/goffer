package scenes

import (
	"image"
)

var (
	Scenes = map[string]Scene{
		"template":    &Template{},
		"fulkonstett": &frameFulkonstOne{},
		"fulkonsttvå": &frameFulkonstTwo{},
		"circ0":       &OnCircle0{},
		"ptbend0":     &PtBend0{},
		"imgimport":   &ImgImport{},
		"lines":       &Lines{},
		"diffeq":      NewDiffEq(),
		"svgepicycle": &SvgEpicycle{},
		"epiopti":     &EpiOpti{},
		"epismooth":   &EpiSmooth{},
		"djanl":       &Djanl{},
	}
)

func LastScene() string {
	return "djanl"
}

type Scene interface {
	Init()
	Frame(t float64) image.Image
}
