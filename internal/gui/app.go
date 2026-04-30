package gui

import (
	"errors"
	"image"
	"log/slog"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"github.com/kristyancarvalho/ak820pro/internal/config"
	"github.com/kristyancarvalho/ak820pro/internal/gui/windows"
	"github.com/kristyancarvalho/ak820pro/internal/hid"
	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
)

type WindowID int

const (
	WindowMode WindowID = iota
	WindowThemeEditor
	WindowTFTScreen
	WindowMacros
)

type AppState struct {
	Dev          hid.Device
	KB           *keyboard.Keyboard
	Cfg          *config.Config
	ActiveWindow WindowID
	Error        string
}

type App struct {
	window        *app.Window
	theme         *material.Theme
	explorer      *explorer.Explorer
	logger        *slog.Logger
	ops           op.Ops
	state         AppState
	startupFailed bool
	exit          widget.Clickable
	modeTab       widget.Clickable
	themeTab      widget.Clickable
	tftTab        widget.Clickable
	macrosTab     widget.Clickable
	modeWindow    *windows.ModeWindow
	themeWindow   *windows.ThemeEditorWindow
	tftWindow     *windows.TFTScreenWindow
	macrosWindow  *windows.MacrosWindow
}

func New(window *app.Window, logger *slog.Logger) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	state := AppState{
		Cfg:          cfg,
		ActiveWindow: WindowMode,
	}
	appValue := &App{
		window:       window,
		theme:        NewTheme(),
		explorer:     explorer.NewExplorer(window),
		logger:       logger,
		state:        state,
		modeWindow:   windows.NewModeWindow(cfg.LastMode),
		themeWindow:  windows.NewThemeEditorWindow(cfg),
		tftWindow:    windows.NewTFTScreenWindow(),
		macrosWindow: windows.NewMacrosWindow(),
	}
	device, err := hid.Open(keyboard.VendorID, keyboard.ProductID)
	if errors.Is(err, hid.ErrDeviceNotFound) {
		device, err = hid.Open(keyboard.LegacyVendorID, keyboard.LegacyProductID)
	}
	if err != nil {
		appValue.startupFailed = true
		if errors.Is(err, hid.ErrDeviceNotFound) {
			appValue.state.Error = "Ajazz AK820 Pro not found. Connect the keyboard in wired mode and reopen the app."
		} else {
			appValue.state.Error = err.Error()
			if appValue.logger != nil {
				appValue.logger.Error("open device", "err", err)
			}
		}
		return appValue, nil
	}
	appValue.state.Dev = device
	appValue.state.KB = keyboard.New(device)
	return appValue, nil
}

func (a *App) Run() error {
	defer a.closeDevice()
	for {
		event := a.window.Event()
		a.explorer.ListenEvents(event)
		switch e := event.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&a.ops, e)
			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) closeDevice() {
	if a.state.Dev != nil {
		_ = a.state.Dev.Close()
	}
}

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	paint.FillShape(gtx.Ops, a.theme.Palette.Bg, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())
	if a.startupFailed {
		for a.exit.Clicked(gtx) {
			a.window.Perform(system.ActionClose)
		}
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(24), Right: unit.Dp(24), Top: unit.Dp(24), Bottom: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(
					gtx,
					layout.Rigid(material.H5(a.theme, "Device Error").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
					}),
					layout.Rigid(material.Body1(a.theme, a.state.Error).Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Height: unit.Dp(18)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Button(a.theme, &a.exit, "Exit").Layout(gtx)
					}),
				)
			})
		})
	}

	for a.modeTab.Clicked(gtx) {
		a.state.ActiveWindow = WindowMode
	}
	for a.themeTab.Clicked(gtx) {
		a.state.ActiveWindow = WindowThemeEditor
	}
	for a.tftTab.Clicked(gtx) {
		a.state.ActiveWindow = WindowTFTScreen
	}
	for a.macrosTab.Clicked(gtx) {
		a.state.ActiveWindow = WindowMacros
	}

	return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Top: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(
			gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.layoutTabs(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return a.layoutActiveWindow(gtx)
			}),
		)
	})
}

func (a *App) layoutTabs(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutTabButton(gtx, &a.modeTab, WindowMode, "Mode")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutTabButton(gtx, &a.themeTab, WindowThemeEditor, "Theme Editor")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutTabButton(gtx, &a.tftTab, WindowTFTScreen, "TFT Screen")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutTabButton(gtx, &a.macrosTab, WindowMacros, "Macros")
		}),
	)
}

func (a *App) layoutTabButton(gtx layout.Context, button *widget.Clickable, current WindowID, label string) layout.Dimensions {
	if a.state.ActiveWindow == current {
		label = "[" + label + "]"
	}
	return material.Button(a.theme, button, label).Layout(gtx)
}

func (a *App) layoutActiveWindow(gtx layout.Context) layout.Dimensions {
	return widget.Border{
		Color:        a.theme.Palette.Fg,
		CornerRadius: unit.Dp(10),
		Width:        unit.Dp(1),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Top: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch a.state.ActiveWindow {
			case WindowThemeEditor:
				return a.themeWindow.Layout(gtx, a.theme, a.state.KB, a.state.Cfg)
			case WindowTFTScreen:
				return a.tftWindow.Layout(gtx, a.theme, a.explorer, a.state.KB, a.window.Invalidate)
			case WindowMacros:
				return a.macrosWindow.Layout(gtx, a.theme, a.state.KB)
			default:
				return a.modeWindow.Layout(gtx, a.theme, a.state.KB, a.state.Cfg)
			}
		})
	})
}
