package keyboard

import (
	"bytes"
	"testing"

	"github.com/kristyancarvalho/ak820pro/internal/hid"
	"github.com/stretchr/testify/require"
)

func TestSetModeWritesExpectedReports(t *testing.T) {
	t.Parallel()
	device := &hid.FakeDevice{}
	board := New(device)

	err := board.SetMode(ModeConfig{
		Mode:       Breath,
		R:          10,
		G:          20,
		B:          30,
		Rainbow:    true,
		Brightness: 40,
		Speed:      50,
		Direction:  Right,
	})
	require.NoError(t, err)
	require.Len(t, device.Writes, 2)
	require.Len(t, device.Writes[0], ReportSize)
	require.Equal(t, []byte{ReportID, CmdSetMode, byte(Breath), 10, 20, 30, 1, 40, 50, byte(Right)}, device.Writes[0][:10])
	require.Equal(t, ReportID, device.Writes[1][0])
	require.Equal(t, CmdCommit, device.Writes[1][1])
}

func TestSetThemeProducesNineChunks(t *testing.T) {
	t.Parallel()
	device := &hid.FakeDevice{}
	board := New(device)

	err := board.SetTheme([][3]uint8{{1, 2, 3}, {4, 5, 6}})
	require.NoError(t, err)
	require.Len(t, device.Writes, NumChunks+1)
	for i := 0; i < NumChunks; i++ {
		require.Len(t, device.Writes[i], ImageChunkSize+2)
		require.Equal(t, ReportID, device.Writes[i][0])
		require.Equal(t, CmdSetPerKey, device.Writes[i][1])
	}
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6}, device.Writes[0][2:8])
	require.Equal(t, byte(0), device.Writes[0][8])
	require.Equal(t, byte(0xFF), device.Writes[0][2+KeyCount*ThemeBytesPerKey])
	require.Equal(t, CmdCommit, device.Writes[len(device.Writes)-1][1])
}

func TestSendTFTImageWritesChunkSequence(t *testing.T) {
	t.Parallel()
	device := &hid.FakeDevice{}
	board := New(device)
	data := bytes.Repeat([]byte{0xAB}, TFTWidth*TFTHeight*TFTBytesPerPixel)

	err := board.SendTFTImage(data)
	require.NoError(t, err)
	require.Len(t, device.Writes, NumChunks+1)
	for i := 0; i < NumChunks; i++ {
		require.Equal(t, ReportID, device.Writes[i][0])
		require.Equal(t, CmdImageChunk, device.Writes[i][1])
	}
	require.Equal(t, CmdImageFinal, device.Writes[len(device.Writes)-1][1])
}
