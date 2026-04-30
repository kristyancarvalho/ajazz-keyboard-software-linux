package widgets

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ColorPicker struct {
	red         widget.Float
	green       widget.Float
	blue        widget.Float
	last        [3]uint8
	initialized bool
}

func NewColorPicker(initial [3]uint8) *ColorPicker {
	picker := &ColorPicker{}
	picker.Set(initial)
	return picker
}

func (p *ColorPicker) Set(value [3]uint8) {
	p.red.Value = float32(value[0]) / 255
	p.green.Value = float32(value[1]) / 255
	p.blue.Value = float32(value[2]) / 255
	p.last = value
	p.initialized = true
}

func (p *ColorPicker) Value() [3]uint8 {
	return [3]uint8{
		floatToByte(p.red.Value),
		floatToByte(p.green.Value),
		floatToByte(p.blue.Value),
	}
}

func (p *ColorPicker) Layout(gtx layout.Context, th *material.Theme, onChange func([3]uint8)) layout.Dimensions {
	if !p.initialized {
		p.Set([3]uint8{})
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutSlider(gtx, th, "Red", &p.red, p.Value()[0])
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutSlider(gtx, th, "Green", &p.green, p.Value()[1])
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutSlider(gtx, th, "Blue", &p.blue, p.Value()[2])
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutSwatch(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			value := p.Value()
			if value != p.last {
				p.last = value
				if onChange != nil {
					onChange(value)
				}
			}
			return layout.Dimensions{}
		}),
	)
}

func (p *ColorPicker) layoutSlider(gtx layout.Context, th *material.Theme, label string, slider *widget.Float, value uint8) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			text := fmt.Sprintf("%s: %d", label, value)
			return material.Body1(th, text).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, slider).Layout(gtx)
		}),
	)
}

func (p *ColorPicker) layoutSwatch(gtx layout.Context, th *material.Theme) layout.Dimensions {
	current := p.Value()
	size := image.Pt(gtx.Dp(unit.Dp(72)), gtx.Dp(unit.Dp(48)))
	return widget.Border{
		Color:        th.Palette.Fg,
		CornerRadius: unit.Dp(8),
		Width:        unit.Dp(1),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(size)
		paint.FillShape(gtx.Ops, toNRGBA(current), clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(unit.Dp(8))).Op(gtx.Ops))
		return layout.Dimensions{Size: size}
	})
}

func floatToByte(value float32) uint8 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return uint8(value * 255)
}

func toNRGBA(value [3]uint8) color.NRGBA {
	return color.NRGBA{R: value[0], G: value[1], B: value[2], A: 0xFF}
}
