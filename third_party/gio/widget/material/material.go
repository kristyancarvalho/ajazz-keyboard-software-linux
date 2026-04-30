package material

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Theme struct {
	TextSize unit.Sp
	Palette  Palette
}

type Palette struct {
	Bg         color.NRGBA
	Fg         color.NRGBA
	ContrastBg color.NRGBA
	ContrastFg color.NRGBA
}

type LabelStyle struct {
	Text  string
	Color color.NRGBA
}

type ButtonStyle struct{}

type RadioButtonStyle struct{}

type CheckBoxStyle struct{}

type EditorStyle struct{}

type SliderStyle struct{}

func NewTheme() *Theme {
	return &Theme{
		Palette: Palette{
			Bg:         color.NRGBA{A: 0xFF},
			Fg:         color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			ContrastBg: color.NRGBA{R: 0x55, G: 0xAA, B: 0xFF, A: 0xFF},
			ContrastFg: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
	}
}

func Body1(th *Theme, text string) LabelStyle {
	return LabelStyle{Text: text, Color: th.Palette.Fg}
}

func Caption(th *Theme, text string) LabelStyle {
	return LabelStyle{Text: text, Color: th.Palette.Fg}
}

func H5(th *Theme, text string) LabelStyle {
	return LabelStyle{Text: text, Color: th.Palette.Fg}
}

func H6(th *Theme, text string) LabelStyle {
	return LabelStyle{Text: text, Color: th.Palette.Fg}
}

func Button(th *Theme, button *widget.Clickable, text string) ButtonStyle {
	return ButtonStyle{}
}

func RadioButton(th *Theme, enum *widget.Enum, key, label string) RadioButtonStyle {
	return RadioButtonStyle{}
}

func CheckBox(th *Theme, state *widget.Bool, label string) CheckBoxStyle {
	return CheckBoxStyle{}
}

func Editor(th *Theme, editor *widget.Editor, hint string) EditorStyle {
	return EditorStyle{}
}

func Slider(th *Theme, value *widget.Float) SliderStyle {
	return SliderStyle{}
}

func (l LabelStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (b ButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (r RadioButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (e EditorStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (s SliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}
