package scenes

import (
	"image"
)

var (
	Scenes = map[string]Project{
		"fulkonstett": &frameFulkonstOne{},
		"fulkonsttvå": &frameFulkonstTwo{},
		"circ0":       &OnCircle0{},
		"ptbend0":     &PtBend0{},
		"imgimport":   &ImgImport{},
	}
)

func LastScene() string {
	return "imgimport"
}

type Project interface {
	Init()
	Frame(t float64) image.Image
}
