package windows

import (
	"sort"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/kristyancarvalho/ak820pro/internal/config"
	"github.com/kristyancarvalho/ak820pro/internal/gui/widgets"
	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
)

type ThemeEditorWindow struct {
	grid        *widgets.KeyGrid
	picker      *widgets.ColorPicker
	name        widget.Editor
	loadChoice  widget.Enum
	applyAll    widget.Clickable
	applyKey    widget.Clickable
	sendTheme   widget.Clickable
	saveTheme   widget.Clickable
	loadTheme   widget.Clickable
	selected    int
	colors      [][3]uint8
	status      string
	initialized bool
}

func NewThemeEditorWindow(cfg *config.Config) *ThemeEditorWindow {
	window := &ThemeEditorWindow{
		grid:   widgets.NewKeyGrid(keyboard.KeyCount),
		picker: widgets.NewColorPicker([3]uint8{255, 255, 255}),
	}
	window.name.SingleLine = true
	window.load(cfg)
	return window
}

func (w *ThemeEditorWindow) Layout(gtx layout.Context, th *material.Theme, kb *keyboard.Keyboard, cfg *config.Config) layout.Dimensions {
	if !w.initialized {
		w.load(cfg)
	}

	for w.applyAll.Clicked(gtx) {
		value := w.picker.Value()
		for i := range w.colors {
			w.colors[i] = value
		}
		w.status = "Applied color to every key"
	}
	for w.applyKey.Clicked(gtx) {
		if w.selected >= 0 && w.selected < len(w.colors) {
			w.colors[w.selected] = w.picker.Value()
			w.status = "Applied color to selected key"
		}
	}
	for w.sendTheme.Clicked(gtx) {
		if kb == nil {
			w.status = "Keyboard not connected"
		} else if err := kb.SetTheme(w.colors); err != nil {
			w.status = err.Error()
		} else {
			w.status = "Theme sent to keyboard"
		}
	}
	for w.saveTheme.Clicked(gtx) {
		name := strings.TrimSpace(w.name.Text())
		if name == "" {
			w.status = "Enter a theme name first"
			continue
		}
		if cfg.Themes == nil {
			cfg.Themes = make(map[string][][3]uint8)
		}
		cfg.Themes[name] = cloneColors(w.colors)
		cfg.LastThemeName = name
		if err := cfg.Save(); err != nil {
			w.status = err.Error()
		} else {
			w.status = "Theme saved"
		}
	}
	for w.loadTheme.Clicked(gtx) {
		name := w.loadChoice.Value
		theme, ok := cfg.Themes[name]
		if !ok {
			w.status = "Select a saved theme"
			continue
		}
		w.colors = normalizeColors(theme)
		if w.selected >= len(w.colors) {
			w.selected = 0
		}
		if len(w.colors) > 0 {
			w.picker.Set(w.colors[w.selected])
		}
		cfg.LastThemeName = name
		if err := cfg.Save(); err != nil {
			w.status = err.Error()
		} else {
			w.status = "Theme loaded"
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(material.H6(th, "Theme Editor").Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Selected key").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.grid.Layout(gtx, th, w.colors, w.selected, func(index int) {
				w.selected = index
				if index >= 0 && index < len(w.colors) {
					w.picker.Set(w.colors[index])
				}
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.picker.Layout(gtx, th, nil)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &w.applyKey, "Apply To Key").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(12)}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &w.applyAll, "Apply To All").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(12)}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &w.sendTheme, "Send Theme").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(18)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Save Theme").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, &w.name, "Theme name").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &w.saveTheme, "Save Theme").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(18)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Saved Themes").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.layoutSavedThemes(gtx, th, cfg)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &w.loadTheme, "Load Selected Theme").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if w.status == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, material.Body1(th, w.status).Layout)
		}),
	)
}

func (w *ThemeEditorWindow) load(cfg *config.Config) {
	loaded := normalizeColors(nil)
	if cfg != nil && cfg.LastThemeName != "" {
		if theme, ok := cfg.Themes[cfg.LastThemeName]; ok {
			loaded = normalizeColors(theme)
			w.loadChoice.Value = cfg.LastThemeName
			w.name.SetText(cfg.LastThemeName)
		}
	}
	w.colors = loaded
	if len(w.colors) > 0 {
		w.picker.Set(w.colors[0])
	}
	w.selected = 0
	w.initialized = true
}

func (w *ThemeEditorWindow) layoutSavedThemes(gtx layout.Context, th *material.Theme, cfg *config.Config) layout.Dimensions {
	names := themeNames(cfg)
	if len(names) == 0 {
		return material.Body1(th, "No saved themes yet").Layout(gtx)
	}
	children := make([]layout.FlexChild, 0, len(names)*2)
	for i, name := range names {
		current := name
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.RadioButton(th, &w.loadChoice, current, current).Layout(gtx)
		}))
		if i != len(names)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(6)}.Layout(gtx)
			}))
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func themeNames(cfg *config.Config) []string {
	if cfg == nil || len(cfg.Themes) == 0 {
		return nil
	}
	names := make([]string, 0, len(cfg.Themes))
	for name := range cfg.Themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func normalizeColors(colors [][3]uint8) [][3]uint8 {
	normalized := make([][3]uint8, keyboard.KeyCount)
	if len(colors) == 0 {
		return normalized
	}
	copy(normalized, colors)
	return normalized
}

func cloneColors(colors [][3]uint8) [][3]uint8 {
	cloned := make([][3]uint8, len(colors))
	copy(cloned, colors)
	return cloned
}
