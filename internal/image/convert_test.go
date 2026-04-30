package image

import (
	stdimage "image"
	"image/color"
	"testing"

	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
	"github.com/stretchr/testify/require"
)

func TestConvertForTFTLengthAndDifference(t *testing.T) {
	t.Parallel()
	red := stdimage.NewNRGBA(stdimage.Rect(0, 0, 4, 4))
	blue := stdimage.NewNRGBA(stdimage.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			red.Set(x, y, color.NRGBA{R: 255, A: 255})
			blue.Set(x, y, color.NRGBA{B: 255, A: 255})
		}
	}

	redData, err := ConvertForTFT(red)
	require.NoError(t, err)
	blueData, err := ConvertForTFT(blue)
	require.NoError(t, err)
	require.Len(t, redData, keyboard.TFTWidth*keyboard.TFTHeight*keyboard.TFTBytesPerPixel)
	require.Len(t, blueData, keyboard.TFTWidth*keyboard.TFTHeight*keyboard.TFTBytesPerPixel)
	require.NotEqual(t, redData, blueData)
}
