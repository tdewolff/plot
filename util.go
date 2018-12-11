package main

import (
	"github.com/tdewolff/canvas"
)

type Projection struct {
	xoffset, yoffset float64
	xscale, yscale   float64
}

func (p Projection) X(x float64) float64 {
	return p.xoffset + p.xscale*x
}

func (p Projection) Y(y float64) float64 {
	return p.yoffset + p.yscale*y
}

type TextAlign int

const (
	AlignLeft TextAlign = iota
	AlignCenter
	AlignRight
	AlignTop
	AlignMiddle
	AlignBottom
)

type Label struct {
	face   canvas.FontFace
	text   string
	x, y   float64
	rot    float64
	halign TextAlign
	valign TextAlign
}

func (l Label) Draw(c canvas.C, proj Projection) {
	x, y := l.x, l.y
	w, h := l.face.BBox(l.text)

	if l.halign == AlignCenter {
		x -= w / 2.0
	} else if l.halign == AlignRight {
		x -= w
	}

	y -= l.face.LineHeight()
	if l.valign == AlignMiddle {
		y += h / 2.0
	} else if l.valign == AlignBottom {
		y -= h
	}

	if l.rot != 0.0 {
		panic("not implemented")
	}

	c.SetFont(l.face)
	c.DrawText(proj.X(x), proj.Y(y), l.text)
}
