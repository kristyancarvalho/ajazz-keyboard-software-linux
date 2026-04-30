package clip

import (
	"image"

	"gioui.org/op"
)

type Shape struct{}

type Rect image.Rectangle

type RRect struct{}

func (r Rect) Op() Shape {
	return Shape{}
}

func UniformRRect(rect image.Rectangle, radius int) RRect {
	return RRect{}
}

func (r RRect) Op(ops *op.Ops) Shape {
	return Shape{}
}
