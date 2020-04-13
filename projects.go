package main

import (
	"image"
)

var (
	Projects = map[string]Project{
		"fulkonstett": &frameFulkonstOne{},
		"fulkonsttvå": &frameFulkonstTwo{},
		"circ0":       &OnCircle0{},
	}

	Pstr = "circ0"
	P    = Projects[Pstr]
)

type Project interface {
	Init()
	Frame(t float64) *image.Paletted
}
