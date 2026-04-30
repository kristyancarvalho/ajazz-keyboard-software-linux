package paint

import (
	"image"
	"image/color"

	"gioui.org/op"
)

type ImageOp struct {
	Src image.Image
}

func NewImageOp(src image.Image) ImageOp {
	return ImageOp{Src: src}
}

func FillShape(ops *op.Ops, color color.NRGBA, shape any) {}
