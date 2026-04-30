package explorer

import (
	"errors"
	"io"

	"gioui.org/app"
)

var (
	ErrUserDecline  = errors.New("user declined selection")
	ErrNotAvailable = errors.New("file picker not available")
)

type Explorer struct{}

func NewExplorer(window *app.Window) *Explorer {
	return &Explorer{}
}

func (e *Explorer) ListenEvents(event any) {}

func (e *Explorer) ChooseFile(pattern string) (io.ReadCloser, error) {
	return nil, ErrNotAvailable
}
