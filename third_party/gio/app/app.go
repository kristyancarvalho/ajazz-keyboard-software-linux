package app

import (
	"image"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type Option func(*Window)

type Window struct {
	size      image.Point
	frameSent bool
	closeSent bool
}

type FrameEvent struct{}

type DestroyEvent struct {
	Err error
}

func NewWindow(options ...Option) *Window {
	window := &Window{
		size: image.Pt(900, 600),
	}
	for _, option := range options {
		if option != nil {
			option(window)
		}
	}
	return window
}

func Title(title string) Option {
	return func(w *Window) {}
}

func Size(width, height unit.Dp) Option {
	return func(w *Window) {
		w.size = image.Pt(int(width), int(height))
	}
}

func (w *Window) Event() any {
	if !w.frameSent {
		w.frameSent = true
		return FrameEvent{}
	}
	if !w.closeSent {
		w.closeSent = true
		return DestroyEvent{}
	}
	return DestroyEvent{}
}

func (w *Window) Perform(action system.Action) {}

func (w *Window) Invalidate() {}

func NewContext(ops *op.Ops, event FrameEvent) layout.Context {
	return layout.Context{
		Ops: ops,
		Constraints: layout.Constraints{
			Max: image.Pt(900, 600),
		},
		Metric: unit.Metric{PxPerDp: 1},
	}
}

func (e FrameEvent) Frame(ops *op.Ops) {}
