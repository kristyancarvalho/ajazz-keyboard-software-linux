package windows

import (
	"errors"
	stdimage "image"
	"image/png"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	akimage "github.com/kristyancarvalho/ak820pro/internal/image"
	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
	xdraw "golang.org/x/image/draw"
)

type tftLoadResult struct {
	preview stdimage.Image
	data    []byte
	err     error
}

type TFTScreenWindow struct {
	chooseFile widget.Clickable
	sendImage  widget.Clickable
	preview    stdimage.Image
	converted  []byte
	status     string
	picking    bool
	results    chan tftLoadResult
}

func NewTFTScreenWindow() *TFTScreenWindow {
	return &TFTScreenWindow{
		results: make(chan tftLoadResult, 1),
	}
}

func (w *TFTScreenWindow) Layout(gtx layout.Context, th *material.Theme, chooser *explorer.Explorer, kb *keyboard.Keyboard, invalidate func()) layout.Dimensions {
	w.consumeResults()

	for w.chooseFile.Clicked(gtx) {
		if w.picking {
			continue
		}
		w.picking = true
		w.status = "Waiting for PNG selection"
		go w.pickFile(chooser, invalidate)
	}
	for w.sendImage.Clicked(gtx) {
		if kb == nil {
			w.status = "Keyboard not connected"
		} else if len(w.converted) == 0 {
			w.status = "Choose a PNG first"
		} else if err := kb.SendTFTImage(w.converted); err != nil {
			w.status = err.Error()
		} else {
			w.status = "Image sent to keyboard"
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(material.H6(th, "TFT Screen").Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &w.chooseFile, "Choose PNG").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(12)}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &w.sendImage, "Send To Keyboard").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if w.preview == nil {
				return material.Body1(th, "No image selected").Layout(gtx)
			}
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				preview := widget.Image{
					Src:   paint.NewImageOp(w.preview),
					Fit:   widget.Contain,
					Scale: 1.0 / gtx.Metric.PxPerDp,
				}
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints = layout.Exact(stdimage.Pt(gtx.Dp(unit.Dp(220)), gtx.Dp(unit.Dp(220))))
					return preview.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if w.status == "" {
				return layout.Dimensions{}
			}
			return material.Body1(th, w.status).Layout(gtx)
		}),
	)
}

func (w *TFTScreenWindow) consumeResults() {
	for {
		select {
		case result := <-w.results:
			w.picking = false
			if result.err != nil {
				w.status = result.err.Error()
				continue
			}
			w.preview = result.preview
			w.converted = result.data
			w.status = "PNG prepared for the TFT screen"
		default:
			return
		}
	}
}

func (w *TFTScreenWindow) pickFile(chooser *explorer.Explorer, invalidate func()) {
	result := tftLoadResult{}
	defer func() {
		w.results <- result
		if invalidate != nil {
			invalidate()
		}
	}()
	if chooser == nil {
		result.err = errors.New("file picker is not available")
		return
	}
	reader, err := chooser.ChooseFile(".png")
	if err != nil {
		if errors.Is(err, explorer.ErrUserDecline) {
			result.err = errors.New("file selection cancelled")
			return
		}
		result.err = err
		return
	}
	defer reader.Close()
	src, err := png.Decode(reader)
	if err != nil {
		result.err = err
		return
	}
	converted, err := akimage.ConvertForTFT(src)
	if err != nil {
		result.err = err
		return
	}
	result.preview = resizePreview(src)
	result.data = converted
}

func resizePreview(src stdimage.Image) stdimage.Image {
	dst := stdimage.NewNRGBA(stdimage.Rect(0, 0, keyboard.TFTWidth, keyboard.TFTHeight))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)
	return dst
}
