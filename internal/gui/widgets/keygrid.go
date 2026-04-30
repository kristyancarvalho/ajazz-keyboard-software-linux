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

type KeyGrid struct {
	keys []widget.Clickable
	rows []int
}

func NewKeyGrid(keyCount int) *KeyGrid {
	rows := []int{16, 16, 15, 14, 11, 9}
	sum := 0
	for _, count := range rows {
		sum += count
	}
	if keyCount > sum {
		rows[len(rows)-1] += keyCount - sum
	}
	grid := &KeyGrid{
		keys: make([]widget.Clickable, keyCount),
		rows: rows,
	}
	return grid
}

func (g *KeyGrid) Layout(gtx layout.Context, th *material.Theme, colors [][3]uint8, selected int, onKeySelected func(int)) layout.Dimensions {
	children := make([]layout.FlexChild, 0, len(g.rows)*2)
	index := 0
	for rowIndex, rowLength := range g.rows {
		start := index
		length := rowLength
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.layoutRow(gtx, th, colors, selected, start, length, onKeySelected)
		}))
		if rowIndex != len(g.rows)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(6)}.Layout(gtx)
			}))
		}
		index += rowLength
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (g *KeyGrid) layoutRow(gtx layout.Context, th *material.Theme, colors [][3]uint8, selected, start, length int, onKeySelected func(int)) layout.Dimensions {
	children := make([]layout.FlexChild, 0, length*2)
	for i := 0; i < length; i++ {
		index := start + i
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.layoutKey(gtx, th, colors, selected, index, onKeySelected)
		}))
		if i != length-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Width: unit.Dp(6)}.Layout(gtx)
			}))
		}
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (g *KeyGrid) layoutKey(gtx layout.Context, th *material.Theme, colors [][3]uint8, selected, index int, onKeySelected func(int)) layout.Dimensions {
	if index >= len(g.keys) {
		return layout.Dimensions{}
	}
	for g.keys[index].Clicked(gtx) {
		if onKeySelected != nil {
			onKeySelected(index)
		}
	}
	size := image.Pt(gtx.Dp(unit.Dp(30)), gtx.Dp(unit.Dp(30)))
	fill := color.NRGBA{A: 0xFF}
	if index < len(colors) {
		fill = toNRGBA(colors[index])
	}
	border := th.Palette.Fg
	if index == selected {
		border = th.Palette.ContrastBg
	}
	return g.keys[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(size)
		return widget.Border{
			Color:        border,
			CornerRadius: unit.Dp(6),
			Width:        unit.Dp(2),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, fill, clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(unit.Dp(6))).Op(gtx.Ops))
			label := material.Caption(th, fmt.Sprintf("%d", index+1))
			label.Color = contrastFor(fill)
			return layout.Center.Layout(gtx, label.Layout)
		})
	})
}

func contrastFor(fill color.NRGBA) color.NRGBA {
	luminance := (int(fill.R)*299 + int(fill.G)*587 + int(fill.B)*114) / 1000
	if luminance < 140 {
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	}
	return color.NRGBA{R: 0x10, G: 0x14, B: 0x1A, A: 0xFF}
}
