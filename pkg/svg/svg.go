package svg

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type OperationType int

const (
	Move OperationType = iota + 1
	MoveRel
	Cubic
	CubicRel // relative to prior instr
	Close    // use the first point
)

type Svg struct {
	Groups []Group `xml:"g"`
}

func ParseSvg(r io.Reader) (*Svg, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	svg := &Svg{}
	err = xml.Unmarshal(buf, svg)
	return svg, err
}

type Group struct {
	Id    string // attribute
	Paths []Path `xml:"path"`
}

type Path struct {
	Id         string
	Operations []Operation
}

func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// need to consume the element?
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	fmt.Printf("Unmarshal Path Local: %v, Space: %v\n", start.Name.Local, start.Name.Space)
	for _, a := range start.Attr {
		name := strings.ToLower(a.Name.Local)
		if name == "id" {
			p.Id = a.Value
		} else if name == "d" {
			ops, err := parsePath(a.Value)
			if err != nil {
				return err
			}

			p.Operations = ops
		}
	}

	return nil
}

type Operation struct {
	Type   OperationType
	Points []Point
}

type Point struct {
	X, Y float64
}

func parsePath(ops string) ([]Operation, error) {
	parser := pathParser{
		tokens:     strings.Fields(ops),
		operations: []Operation{},
	}

	return parser.Parse()
}

type lexFunc func(string) lexFunc
type pathParser struct {
	at         int
	tokens     []string
	operations []Operation
	curveInit  bool
	relative   bool
	err        error
}

func (p *pathParser) Parse() ([]Operation, error) {
	if len(p.tokens) < 2 {
		return nil, fmt.Errorf("too few tokens")
	}

	var lex lexFunc = p.operation
	for lex != nil && p.at < len(p.tokens) {
		lex = lex(p.tokens[p.at])
		p.next()
	}

	return p.operations, p.err
}

func (p *pathParser) next() bool {
	p.at++
	return p.at < len(p.tokens)
}

func (p *pathParser) get() string {
	if p.at < len(p.tokens) {
		return p.tokens[p.at]
	} else {
		return ""
	}
}

func (p *pathParser) peek() string {
	next := p.at + 1
	if next < len(p.tokens) {
		return p.tokens[next]
	} else {
		return ""
	}
}

func isLetter(s string) bool {

	/*
	   M = moveto
	   L = lineto
	   H = horizontal lineto
	   V = vertical lineto
	   C = curveto
	   S = smooth curveto
	   Q = quadratic Bézier curve
	   T = smooth quadratic Bézier curveto
	   A = elliptical Arc
	   Z = closepath
	*/

	switch strings.ToLower(s) {
	case "m":
		return true
	case "c":
		return true
	case "z":
		return true
	default:
		return false
	}
}

func parsePoint(s string) (Point, error) {
	splt := strings.Split(s, ",")

	p := Point{}
	if len(splt) != 2 {
		return p, fmt.Errorf("bad point string: %v", s)
	}

	x, err := strconv.ParseFloat(splt[0], 64)
	if err != nil {
		return p, err
	}
	y, err := strconv.ParseFloat(splt[1], 64)
	if err != nil {
		return p, err
	}

	p.X, p.Y = x, y
	return p, nil
}

func (p *pathParser) operation(s string) lexFunc {
	// empty is valid; it is the eof indicator
	if s == "" {
		return nil
	}

	if !isLetter(s) {
		return p.errorf("bad draw instruction: %v", s)
	}

	r, _ := utf8.DecodeRuneInString(s)
	if unicode.IsUpper(r) {
		p.relative = false
	} else {
		p.relative = true
	}

	l := strings.ToLower(s)
	if l == "m" {
		return p.move
	} else if l == "c" {
		return p.cubic
	} else if l == "z" {
		return p.close(s) // evaluate, because there are no params
	}

	return p.errorf("internal operation error")
}

func (p *pathParser) move(s string) lexFunc {
	pt, err := parsePoint(s)
	if err != nil {
		return p.errorf("move operation requires point")
	}

	p.operations = append(p.operations, Operation{
		Type:   Move,
		Points: []Point{pt},
	})

	if p.relative {
		p.operations[len(p.operations)-1].Type = MoveRel
	}

	return p.operation
}

func (p *pathParser) cubic(s string) lexFunc {
	pt, err := parsePoint(s)
	if err != nil {
		return p.errorf("curve operation requires point")
	}

	if !p.curveInit {
		p.operations = append(p.operations, Operation{
			Type:   Cubic,
			Points: []Point{pt},
		})

		if p.relative {
			p.operations[len(p.operations)-1].Type = CubicRel
		}

		p.curveInit = true
	} else {
		l := len(p.operations) - 1
		p.operations[l].Points = append(p.operations[l].Points, pt)
	}

	if isLetter(p.peek()) {
		p.curveInit = false
		return p.operation
	}

	return p.cubic
}

func (p *pathParser) close(s string) lexFunc {
	p.operations = append(p.operations, Operation{
		Type: Close,
	})

	return nil // always exit here
}

func (p *pathParser) errorf(s string, args ...string) lexFunc {
	p.err = fmt.Errorf(s, args)
	return nil
}
