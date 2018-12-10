package main

import "github.com/tdewolff/canvas"

type Plot struct {
	lines []*Line

	Margin float64
}

func New() *Plot {
	return &Plot{
		lines: []*Line{},

		Margin: 10.0,
	}
}

func (p *Plot) Add(l *Line) {
	p.lines = append(p.lines, l)
}

func (p *Plot) Draw(c canvas.C, w, h float64) {
	c.Open(w+2*p.Margin, h+2*p.Margin)

	rect := canvas.Rect{p.Margin, p.Margin, w, h}

	axes := canvas.Rectangle(rect.X, rect.Y, rect.W, rect.H)
	axes = axes.Stroke(0.3, canvas.RoundCapper, canvas.RoundJoiner, 1.0)
	c.DrawPath(0.0, 0.0, axes)
	for _, l := range p.lines {
		l.Draw(c, rect)
	}
}

type Line struct {
	xs, ys []float64
}

func NewLine(xs, ys []float64) *Line {
	if len(xs) != len(ys) {
		panic("number of x and y data points do not match")
	}

	return &Line{
		xs: xs,
		ys: ys,
	}
}

func (l *Line) Draw(c canvas.C, rect canvas.Rect) {
	p := &canvas.Path{}
	if len(l.xs) > 0 {
		p.MoveTo(l.xs[0], rect.H-l.ys[0])
	}
	for i := 1; i < len(l.xs); i++ {
		p.LineTo(l.xs[i], rect.H-l.ys[i])
	}
	p = p.Stroke(0.3, canvas.RoundCapper, canvas.RoundJoiner, 1.0)
	c.DrawPath(rect.X, rect.Y, p)
}
