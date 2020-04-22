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
	}
)

func LastScene() string {
	return "diffeq"
}

type Scene interface {
	Init()
	Frame(t float64) image.Image
}
