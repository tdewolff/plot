package main

import (
	"fmt"
	"math"

	"github.com/tdewolff/canvas"
)

type AxisScale int

const (
	LinearScale AxisScale = iota
	LogScale
)

type Position int

const (
	Left Position = iota
	Right
	Top
	Bottom
)

type Range struct {
	Min, Max float64
}

func (l Range) Merge(r Range) Range {
	return Range{
		Min: math.Min(l.Min, r.Min),
		Max: math.Max(l.Max, r.Max),
	}
}

type Axis struct {
	Position
	Range
	Scale        AxisScale
	Ticks        []float64
	LabelPadding float64
}

func (a Axis) LabelSpace(c canvas.C, font canvas.Font) float64 {
	face := font.Face(3.0)
	space := 0.0
	for _, pos := range a.Ticks {
		t := fmt.Sprintf("%g", pos)
		tw, th := face.BBox(t)
		switch a.Position {
		case Left:
			space = math.Max(space, tw)
		case Bottom:
			space = math.Max(space, th)
		default:
			panic("not implemented")
		}
	}
	return space + a.LabelPadding
}

func (a Axis) Draw(c canvas.C, proj Projection, axes Axes, font canvas.Font) {
	face := font.Face(3.0)
	p := &canvas.Path{}
	for _, pos := range a.Ticks {
		switch a.Position {
		case Left:
			label := Label{face, fmt.Sprintf("%g", pos), -a.LabelPadding, pos, 0.0, AlignRight, AlignMiddle}
			label.Draw(c, proj)
			if pos == a.Min || pos == a.Max {
				continue
			}
			p.MoveTo(0.0, pos)
			p.LineTo(2.0, pos)
		case Bottom:
			label := Label{face, fmt.Sprintf("%g", pos), pos, -a.LabelPadding, 0.0, AlignCenter, AlignTop}
			label.Draw(c, proj)
			if pos == a.Min || pos == a.Max {
				continue
			}
			p.MoveTo(pos, 0.0)
			p.LineTo(pos, 2.0)
		default:
			panic("not implemented")
		}
	}
	p = p.Scale(proj.xscale, proj.yscale)
	p = p.Stroke(0.3, canvas.RoundCapper, canvas.RoundJoiner, 1.0)
	c.DrawPath(proj.xoffset, proj.yoffset, p)
}

type Axes struct {
	X, Y Axis
}

func NewAxes(xrange, yrange Range) Axes {
	a := Axes{
		X: Axis{
			Position:     Bottom,
			Range:        xrange,
			Scale:        LinearScale,
			LabelPadding: 1.0,
		},
		Y: Axis{
			Position:     Left,
			Range:        yrange,
			Scale:        LinearScale,
			LabelPadding: 1.0,
		},
	}

	a.X.Ticks, _, _, _ = talbotLinHanrahan(xrange.Min, xrange.Max, 5, free, nil, nil, nil)
	a.Y.Ticks, _, _, _ = talbotLinHanrahan(yrange.Min, yrange.Max, 5, free, nil, nil, nil)
	return a
}

func (a Axes) Draw(c canvas.C, proj Projection, rect canvas.Rect, font canvas.Font) {
	axes := canvas.Rectangle(rect.X, rect.Y, rect.W, rect.H)
	axes = axes.Stroke(0.3, canvas.RoundCapper, canvas.RoundJoiner, 1.0)
	c.DrawPath(0.0, 0.0, axes)

	a.X.Draw(c, proj, a, font)
	a.Y.Draw(c, proj, a, font)
}

type Plot struct {
	title  string
	xlabel string
	ylabel string
	lines  []*Line

	Margin       float64
	TitlePadding float64
	LabelPadding float64
}

func New(title string) *Plot {
	return &Plot{
		title: title,
		lines: []*Line{},

		Margin:       1.0,
		TitlePadding: 5.0,
		LabelPadding: 3.0,
	}
}

func (p *Plot) SetXLabel(xlabel string) {
	p.xlabel = xlabel
}

func (p *Plot) SetYLabel(ylabel string) {
	p.ylabel = ylabel
}

func (p *Plot) Add(l *Line) {
	p.lines = append(p.lines, l)
}

func (p *Plot) Draw(c canvas.C, font canvas.Font, w, h float64) {
	c.Open(w, h)

	titleFace := font.Face(18.0)
	labelFace := font.Face(12.0)

	var xrange, yrange Range
	for _, l := range p.lines {
		lxrange, lyrange := l.Ranges()
		lyrange.Min = 0.0
		lyrange.Max *= 1.10
		xrange = xrange.Merge(lxrange)
		yrange = yrange.Merge(lyrange)
	}
	axes := NewAxes(xrange, yrange)

	topMargin := p.Margin
	if p.title != "" {
		topMargin += titleFace.Metrics().LineHeight + p.TitlePadding
	}
	bottomMargin := p.Margin + axes.X.LabelSpace(c, font)
	if p.ylabel != "" {
		bottomMargin += labelFace.Metrics().LineHeight + p.LabelPadding
	}
	leftMargin := p.Margin + axes.Y.LabelSpace(c, font)
	if p.xlabel != "" {
		leftMargin += labelFace.Metrics().LineHeight + p.LabelPadding
	}
	rightMargin := p.Margin
	rect := canvas.Rect{leftMargin, h - bottomMargin, w - leftMargin - rightMargin, -(h - topMargin - bottomMargin)}

	proj := Projection{
		xoffset: rect.X,
		yoffset: rect.Y,
		xscale:  rect.W / (axes.X.Max - axes.X.Min),
		yscale:  rect.H / (axes.Y.Max - axes.Y.Min),
	}

	axes.Draw(c, proj, rect, font)
	for _, l := range p.lines {
		l.Draw(c, proj)
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

func (l *Line) Ranges() (Range, Range) {
	xmin, xmax := math.Inf(1), math.Inf(-1)
	ymin, ymax := math.Inf(1), math.Inf(-1)
	for i := 0; i < len(l.xs); i++ {
		if l.xs[i] < xmin {
			xmin = l.xs[i]
		}
		if l.xs[i] > xmax {
			xmax = l.xs[i]
		}
		if l.ys[i] < ymin {
			ymin = l.ys[i]
		}
		if l.ys[i] > ymax {
			ymax = l.ys[i]
		}
	}
	return Range{xmin, xmax}, Range{ymin, ymax}
}

func (l *Line) Draw(c canvas.C, proj Projection) {
	p := &canvas.Path{}
	if len(l.xs) > 0 {
		p.MoveTo(l.xs[0], l.ys[0])
	}
	for i := 1; i < len(l.xs); i++ {
		p.LineTo(l.xs[i], l.ys[i])
	}
	p = p.Scale(proj.xscale, proj.yscale)
	p = p.Stroke(0.3, canvas.RoundCapper, canvas.RoundJoiner, 1.0)
	c.DrawPath(proj.xoffset, proj.yoffset, p)
}
