package layout

import (
	"image"

	"gioui.org/op"
	"gioui.org/unit"
)

type Widget func(Context) Dimensions

type Context struct {
	Ops         *op.Ops
	Constraints Constraints
	Metric      unit.Metric
}

type Constraints struct {
	Min image.Point
	Max image.Point
}

type Dimensions struct {
	Size     image.Point
	Baseline int
}

type Axis uint8

type Flex struct {
	Axis Axis
}

type FlexChild struct {
	widget Widget
}

type Inset struct {
	Top    unit.Dp
	Bottom unit.Dp
	Left   unit.Dp
	Right  unit.Dp
}

type Spacer struct {
	Width  unit.Dp
	Height unit.Dp
}

type centerLayout struct{}

const (
	Horizontal Axis = iota
	Vertical
)

var Center centerLayout

func Exact(size image.Point) Constraints {
	return Constraints{Min: size, Max: size}
}

func Rigid(w Widget) FlexChild {
	return FlexChild{widget: w}
}

func Flexed(weight float32, w Widget) FlexChild {
	return FlexChild{widget: w}
}

func (c Context) Dp(value unit.Dp) int {
	scale := c.Metric.PxPerDp
	if scale == 0 {
		scale = 1
	}
	return int(float32(value)*scale + 0.5)
}

func (f Flex) Layout(gtx Context, children ...FlexChild) Dimensions {
	var size image.Point
	for _, child := range children {
		if child.widget == nil {
			continue
		}
		dims := child.widget(gtx)
		if f.Axis == Horizontal {
			size.X += dims.Size.X
			if dims.Size.Y > size.Y {
				size.Y = dims.Size.Y
			}
			continue
		}
		size.Y += dims.Size.Y
		if dims.Size.X > size.X {
			size.X = dims.Size.X
		}
	}
	return Dimensions{Size: size}
}

func (i Inset) Layout(gtx Context, w Widget) Dimensions {
	if w == nil {
		return Dimensions{}
	}
	return w(gtx)
}

func (s Spacer) Layout(gtx Context) Dimensions {
	return Dimensions{Size: image.Pt(gtx.Dp(s.Width), gtx.Dp(s.Height))}
}

func (centerLayout) Layout(gtx Context, w Widget) Dimensions {
	if w == nil {
		return Dimensions{}
	}
	return w(gtx)
}
