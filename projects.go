package main

import (
	"image"
)

var (
	Projects = map[string]Project{
		"fulkonstett": &frameFulkonstOne{},
		"fulkonsttvå": &frameFulkonstTwo{},
		"circ0":       &OnCircle0{},
		"ptbend0":     &PtBend0{},
	}

	ActiveProject = "ptbend0"
	P             = Projects[ActiveProject]
)

type Project interface {
	Init()
	Frame(t float64) image.Image
}
