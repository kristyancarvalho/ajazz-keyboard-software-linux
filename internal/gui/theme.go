package gui

import (
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget/material"
)

func NewTheme() *material.Theme {
	th := material.NewTheme()
	th.TextSize = unit.Sp(14)
	th.Palette = material.Palette{
		Bg:         color.NRGBA{R: 0x10, G: 0x14, B: 0x1A, A: 0xFF},
		Fg:         color.NRGBA{R: 0xF3, G: 0xF6, B: 0xFA, A: 0xFF},
		ContrastBg: color.NRGBA{R: 0x0F, G: 0x8B, B: 0x8D, A: 0xFF},
		ContrastFg: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}
	return th
}
