package windows

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
)

type macroEntry struct {
	editor widget.Editor
	apply  widget.Clickable
}

type MacrosWindow struct {
	slots       []macroEntry
	status      string
	initialized bool
}

func NewMacrosWindow() *MacrosWindow {
	return &MacrosWindow{}
}

func (w *MacrosWindow) Layout(gtx layout.Context, th *material.Theme, kb *keyboard.Keyboard) layout.Dimensions {
	if !w.initialized {
		w.slots = make([]macroEntry, 4)
		for i := range w.slots {
			w.slots[i].editor.SingleLine = true
		}
		w.initialized = true
	}
	children := make([]layout.FlexChild, 0, len(w.slots)*2+3)
	children = append(children, layout.Rigid(material.H6(th, "Macros").Layout))
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
	}))
	for index := range w.slots {
		slotIndex := index
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutSlot(gtx, th, kb, slotIndex)
		}))
		if index != len(w.slots)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
			}))
		}
	}
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		if w.status == "" {
			return layout.Dimensions{}
		}
		return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, material.Body1(th, w.status).Layout)
	}))
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (w *MacrosWindow) layoutSlot(gtx layout.Context, th *material.Theme, kb *keyboard.Keyboard, slotIndex int) layout.Dimensions {
	slot := &w.slots[slotIndex]
	for slot.apply.Clicked(gtx) {
		sequence, err := parseMacroSequence(slot.editor.Text())
		if err != nil {
			w.status = err.Error()
			continue
		}
		if kb == nil {
			w.status = "Keyboard not connected"
			continue
		}
		if err := kb.SetMacro(keyboard.MacroSlot{Index: slotIndex, Sequence: sequence}); err != nil {
			w.status = err.Error()
		} else {
			w.status = fmt.Sprintf("Macro %d applied", slotIndex+1)
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, fmt.Sprintf("Slot %d", slotIndex+1)).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(6)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, &slot.editor, "Bytes like 0x1A 0x2B 44").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &slot.apply, "Apply Macro").Layout(gtx)
			})
		}),
	)
}

func parseMacroSequence(input string) ([]uint8, error) {
	fields := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == ',' || r == ';' || r == '\n' || r == '\t'
	})
	sequence := make([]uint8, 0, len(fields))
	for _, field := range fields {
		base := 10
		value := field
		if strings.HasPrefix(strings.ToLower(field), "0x") {
			base = 16
			value = field[2:]
		}
		parsed, err := strconv.ParseUint(value, base, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid macro byte %q", field)
		}
		sequence = append(sequence, uint8(parsed))
	}
	return sequence, nil
}
