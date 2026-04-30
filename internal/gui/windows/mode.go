package windows

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/kristyancarvalho/ak820pro/internal/config"
	"github.com/kristyancarvalho/ak820pro/internal/gui/widgets"
	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
)

type ModeWindow struct {
	mode        widget.Enum
	direction   widget.Enum
	rainbow     widget.Bool
	brightness  widget.Float
	speed       widget.Float
	apply       widget.Clickable
	picker      *widgets.ColorPicker
	status      string
	initialized bool
}

func NewModeWindow(cfg keyboard.ModeConfig) *ModeWindow {
	window := &ModeWindow{
		picker: widgets.NewColorPicker([3]uint8{cfg.R, cfg.G, cfg.B}),
	}
	window.load(cfg)
	return window
}

func (w *ModeWindow) Layout(gtx layout.Context, th *material.Theme, kb *keyboard.Keyboard, cfg *config.Config) layout.Dimensions {
	if !w.initialized {
		w.load(cfg.LastMode)
	}
	for w.apply.Clicked(gtx) {
		mode := parseMode(w.mode.Value)
		direction := parseDirection(w.direction.Value)
		configValue := keyboard.ModeConfig{
			Mode:       mode,
			R:          w.picker.Value()[0],
			G:          w.picker.Value()[1],
			B:          w.picker.Value()[2],
			Rainbow:    w.rainbow.Value,
			Brightness: sliderByte(w.brightness.Value),
			Speed:      sliderByte(w.speed.Value),
			Direction:  direction,
		}
		if kb == nil {
			w.status = "Keyboard not connected"
		} else if err := kb.SetMode(configValue); err != nil {
			w.status = err.Error()
		} else {
			cfg.LastMode = configValue
			if err := cfg.Save(); err != nil {
				w.status = err.Error()
			} else {
				w.status = "Lighting mode applied"
			}
		}
	}

	selectedMode := parseMode(w.mode.Value)
	allowedDirections := keyboard.DirectionsForMode(selectedMode)
	if !containsDirection(allowedDirections, parseDirection(w.direction.Value)) {
		w.direction.Value = directionKey(allowedDirections[0])
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(material.H6(th, "Lighting Mode").Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutModes(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.CheckBox(th, &w.rainbow, "Rainbow").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if w.rainbow.Value {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return w.picker.Layout(gtx, th, nil)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutSliderSection(gtx, th, "Brightness", &w.brightness, sliderByte(w.brightness.Value))
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutSliderSection(gtx, th, "Speed", &w.speed, sliderByte(w.speed.Value))
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Direction").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutDirections(gtx, th, allowedDirections)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &w.apply, "Apply").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if w.status == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, material.Body1(th, w.status).Layout)
		}),
	)
}

func (w *ModeWindow) load(cfg keyboard.ModeConfig) {
	mode := cfg.Mode
	if mode.String() == "Unknown" {
		mode = keyboard.STATIC
	}
	w.mode.Value = modeKey(mode)
	w.rainbow.Value = cfg.Rainbow
	w.brightness.Value = float32(cfg.Brightness) / 255
	w.speed.Value = float32(cfg.Speed) / 255
	direction := cfg.Direction
	if !containsDirection(keyboard.DirectionsForMode(mode), direction) {
		direction = keyboard.DefaultDirection(mode)
	}
	w.direction.Value = directionKey(direction)
	w.picker.Set([3]uint8{cfg.R, cfg.G, cfg.B})
	w.initialized = true
}

func (w *ModeWindow) layoutModes(gtx layout.Context, th *material.Theme) layout.Dimensions {
	modes := keyboard.LightingMode(0).All()
	children := make([]layout.FlexChild, 0, len(modes)/4+1)
	for start := 0; start < len(modes); start += 4 {
		offset := start
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			end := offset + 4
			if end > len(modes) {
				end = len(modes)
			}
			row := make([]layout.FlexChild, 0, (end-offset)*2)
			for i := offset; i < end; i++ {
				mode := modes[i]
				row = append(row, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &w.mode, modeKey(mode), mode.String()).Layout(gtx)
				}))
				if i != end-1 {
					row = append(row, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(12)}.Layout(gtx)
					}))
				}
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, row...)
		}))
		if start+4 < len(modes) {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(6)}.Layout(gtx)
			}))
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (w *ModeWindow) layoutDirections(gtx layout.Context, th *material.Theme, directions []keyboard.Direction) layout.Dimensions {
	children := make([]layout.FlexChild, 0, len(directions)*2)
	for i, direction := range directions {
		dir := direction
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.RadioButton(th, &w.direction, directionKey(dir), dir.String()).Layout(gtx)
		}))
		if i != len(directions)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Width: unit.Dp(12)}.Layout(gtx)
			}))
		}
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (w *ModeWindow) layoutSliderSection(gtx layout.Context, th *material.Theme, label string, state *widget.Float, value uint8) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, fmt.Sprintf("%s: %d", label, value)).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, state).Layout(gtx)
		}),
	)
}

func sliderByte(value float32) uint8 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return uint8(value * 255)
}

func modeKey(mode keyboard.LightingMode) string {
	return fmt.Sprintf("%d", mode)
}

func directionKey(direction keyboard.Direction) string {
	return fmt.Sprintf("%d", direction)
}

func parseMode(value string) keyboard.LightingMode {
	for _, mode := range keyboard.LightingMode(0).All() {
		if modeKey(mode) == value {
			return mode
		}
	}
	return keyboard.STATIC
}

func parseDirection(value string) keyboard.Direction {
	for _, direction := range []keyboard.Direction{keyboard.LEFT, keyboard.DOWN, keyboard.UP, keyboard.RIGHT} {
		if directionKey(direction) == value {
			return direction
		}
	}
	return keyboard.LEFT
}

func containsDirection(directions []keyboard.Direction, target keyboard.Direction) bool {
	for _, direction := range directions {
		if direction == target {
			return true
		}
	}
	return false
}
