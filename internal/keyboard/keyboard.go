package keyboard

import (
	"errors"
	"io"

	"github.com/kristyancarvalho/ak820pro/internal/hid"
)

type ModeConfig struct {
	Mode       LightingMode `json:"mode"`
	R          uint8        `json:"r"`
	G          uint8        `json:"g"`
	B          uint8        `json:"b"`
	Rainbow    bool         `json:"rainbow"`
	Brightness uint8        `json:"brightness"`
	Speed      uint8        `json:"speed"`
	Direction  Direction    `json:"direction"`
}

type MacroSlot struct {
	Index    int
	Sequence []uint8
}

type Keyboard struct {
	dev hid.Device
}

func New(dev hid.Device) *Keyboard {
	return &Keyboard{dev: dev}
}

func (k *Keyboard) SetMode(cfg ModeConfig) error {
	payload := []byte{
		byte(cfg.Mode),
		cfg.R,
		cfg.G,
		cfg.B,
		boolByte(cfg.Rainbow),
		cfg.Brightness,
		cfg.Speed,
		byte(cfg.Direction),
	}
	if err := k.writeCommand(CmdSetMode, payload); err != nil {
		return err
	}
	return k.Commit()
}

func (k *Keyboard) SetTheme(keyColors [][3]uint8) error {
	data := make([]byte, 0, KeyCount*ThemeBytesPerKey)
	for i := 0; i < KeyCount; i++ {
		if i < len(keyColors) {
			color := keyColors[i]
			data = append(data, color[0], color[1], color[2])
			continue
		}
		data = append(data, 0, 0, 0)
	}
	chunks, err := splitData(data)
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		if err := k.writeChunk(CmdSetPerKey, chunk); err != nil {
			return err
		}
	}
	return k.Commit()
}

func (k *Keyboard) SendTFTImage(data []byte) error {
	chunks, err := splitData(data)
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		if err := k.writeChunk(CmdImageChunk, chunk); err != nil {
			return err
		}
	}
	return k.writeCommand(CmdImageFinal, nil)
}

func (k *Keyboard) SetMacro(slot MacroSlot) error {
	if slot.Index < 0 || slot.Index > 255 {
		return errors.New("macro slot index out of range")
	}
	if len(slot.Sequence) > MaxMacroSequence {
		return errors.New("macro sequence too long")
	}
	payload := make([]byte, 0, 2+len(slot.Sequence))
	payload = append(payload, byte(slot.Index), byte(len(slot.Sequence)))
	payload = append(payload, slot.Sequence...)
	if err := k.writeCommand(CmdSetMacro, payload); err != nil {
		return err
	}
	return k.Commit()
}

func (k *Keyboard) Commit() error {
	return k.writeCommand(CmdCommit, nil)
}

func splitData(data []byte) ([][]byte, error) {
	if len(data) > NumChunks*ImageChunkSize {
		return nil, errors.New("data exceeds transport capacity")
	}
	chunks := make([][]byte, NumChunks)
	for i := range chunks {
		chunks[i] = make([]byte, ImageChunkSize)
		for j := range chunks[i] {
			chunks[i][j] = 0xFF
		}
	}
	offset := 0
	remaining := len(data)
	for i := 0; i < NumChunks && remaining > 0; i++ {
		copySize := ImageChunkSize
		if remaining < copySize {
			copySize = remaining
		}
		copy(chunks[i], data[offset:offset+copySize])
		offset += copySize
		remaining -= copySize
	}
	return chunks, nil
}

func (k *Keyboard) writeCommand(cmd byte, payload []byte) error {
	size := ReportSize
	if len(payload)+2 > size {
		size = len(payload) + 2
	}
	report := make([]byte, size)
	report[0] = ReportID
	report[1] = cmd
	copy(report[2:], payload)
	return k.write(report)
}

func (k *Keyboard) writeChunk(cmd byte, payload []byte) error {
	report := make([]byte, len(payload)+2)
	report[0] = ReportID
	report[1] = cmd
	copy(report[2:], payload)
	return k.write(report)
}

func (k *Keyboard) write(report []byte) error {
	if k == nil || k.dev == nil {
		return errors.New("keyboard device is not available")
	}
	n, err := k.dev.Write(report)
	if err != nil {
		return err
	}
	if n != len(report) {
		return io.ErrShortWrite
	}
	return nil
}

func boolByte(value bool) byte {
	if value {
		return 1
	}
	return 0
}
