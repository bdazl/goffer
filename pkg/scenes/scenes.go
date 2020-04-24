package scenes

import (
	"image"
)

var (
	Scenes = map[string]Scene{
		"fulkonstett": &frameFulkonstOne{},
		"fulkonsttvå": &frameFulkonstTwo{},
		"circ0":       &OnCircle0{},
		"ptbend0":     &PtBend0{},
		"imgimport":   &ImgImport{},
		"lines":       &Lines{},
		"diffeq":      NewDiffEq(),
		"svgepicycle": &SvgEpicycle{},
	}
)

func LastScene() string {
	return "svgepicycle"
}

type Scene interface {
	Init()
	Frame(t float64) image.Image
}
