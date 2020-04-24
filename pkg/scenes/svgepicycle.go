package scenes

import (
	"fmt"
	"image"
	"os"

	"github.com/rustyoz/svg"

	jimage "github.com/HexHacks/goffer/pkg/image"
)

type SvgEpicycle struct {
}

func (se *SvgEpicycle) Init() {
	fil, err := os.Open("assets/gubb.svg")
	panicOn(err)
	defer fil.Close()

	s, err := svg.ParseSvgFromReader(fil, "gubbe", 1.0)
	panicOn(err)

	pis := parseDrawingInstructions(s)
	for _, p := range pis {
		switch p.Kind {
		case svg.MoveInstruction:
			fmt.Println("Move")
		case svg.CircleInstruction:
			fmt.Println("Circle")
		case svg.CurveInstruction:
			fmt.Println("Curve")
		case svg.LineInstruction:
			fmt.Println("Line")
		case svg.HLineInstruction:
			fmt.Println("HLine")
		case svg.CloseInstruction:
			fmt.Println("Close")
		case svg.PaintInstruction:
			fmt.Println("Paint")
		}
	}
}

func parseDrawingInstructions(s *svg.Svg) []*svg.DrawingInstruction {
	di := make([]*svg.DrawingInstruction, 0)
	for _, e := range s.Elements {
		ch, errch := e.ParseDrawingInstructions()
		for {
			select {
			case pi, ok := <-ch:
				fmt.Println("GSDFLKJ")
				if !ok {
					break
				}
				fmt.Println("wakawak")

				di = append(di, pi)
			case err := <-errch:
				panicOn(err)
			}
		}
	}
	return di
}

func (se *SvgEpicycle) Frame(s float64) image.Image {
	img, _ := jimage.New()

	return img
}
