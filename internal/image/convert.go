package image

import (
	"errors"
	"image"

	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
	"golang.org/x/image/draw"
)

func ConvertForTFT(src image.Image) ([]byte, error) {
	if src == nil {
		return nil, errors.New("source image is nil")
	}
	if src.Bounds().Empty() {
		return nil, errors.New("source image has empty bounds")
	}
	dst := image.NewNRGBA(image.Rect(0, 0, keyboard.TFTWidth, keyboard.TFTHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	data := make([]byte, keyboard.TFTWidth*keyboard.TFTHeight*keyboard.TFTBytesPerPixel)
	offset := 0
	for y := 0; y < keyboard.TFTHeight; y++ {
		for x := 0; x < keyboard.TFTWidth; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()
			value := rgb565(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			data[offset] = byte(value >> 8)
			data[offset+1] = byte(value)
			offset += keyboard.TFTBytesPerPixel
		}
	}
	return data, nil
}

func rgb565(r, g, b uint8) uint16 {
	return uint16(r&0xF8)<<8 | uint16(g&0xFC)<<3 | uint16(b>>3)
}
