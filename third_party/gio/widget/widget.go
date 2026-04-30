package widget

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Float struct {
	Value float32
}

type Bool struct {
	Value bool
}

type Enum struct {
	Value string
}

type Clickable struct{}

type Editor struct {
	SingleLine bool
	text       string
}

type Fit uint8

type Image struct {
	Src   paint.ImageOp
	Fit   Fit
	Scale float32
}

type Border struct {
	Color        color.NRGBA
	CornerRadius unit.Dp
	Width        unit.Dp
}

const (
	Contain Fit = iota
)

func (c *Clickable) Clicked(gtx layout.Context) bool {
	return false
}

func (c *Clickable) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	if w == nil {
		return layout.Dimensions{}
	}
	return w(gtx)
}

func (e *Editor) Text() string {
	return e.text
}

func (e *Editor) SetText(text string) {
	e.text = text
}

func (i Image) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (b Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	if w == nil {
		return layout.Dimensions{}
	}
	return w(gtx)
}
