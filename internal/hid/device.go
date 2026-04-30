package hid

import (
	"errors"
	"io"

	sshid "github.com/sstallion/go-hid"
)

var ErrDeviceNotFound = errors.New("ajazz ak820 pro device not found")

type Device interface {
	Write(data []byte) (int, error)
	Read(data []byte) (int, error)
	Close() error
}

type FakeDevice struct {
	Writes     [][]byte
	Reads      [][]byte
	WriteErr   error
	ReadErr    error
	CloseErr   error
	ReadOffset int
}

type backend interface {
	Enumerate(vendorID, productID uint16, visit func(info *sshid.DeviceInfo) error) error
	OpenPath(path string) (Device, error)
}

type packageBackend struct{}

type wrappedDevice struct {
	device *sshid.Device
}

func Open(vendorID, productID uint16) (Device, error) {
	return openWithBackend(packageBackend{}, vendorID, productID)
}

func (d *FakeDevice) Write(data []byte) (int, error) {
	if d.WriteErr != nil {
		return 0, d.WriteErr
	}
	copied := append([]byte(nil), data...)
	d.Writes = append(d.Writes, copied)
	return len(data), nil
}

func (d *FakeDevice) Read(data []byte) (int, error) {
	if d.ReadErr != nil {
		return 0, d.ReadErr
	}
	if d.ReadOffset >= len(d.Reads) {
		return 0, io.EOF
	}
	response := d.Reads[d.ReadOffset]
	d.ReadOffset++
	n := copy(data, response)
	return n, nil
}

func (d *FakeDevice) Close() error {
	return d.CloseErr
}

func (b packageBackend) Enumerate(vendorID, productID uint16, visit func(info *sshid.DeviceInfo) error) error {
	if err := sshid.Init(); err != nil {
		return err
	}
	return sshid.Enumerate(vendorID, productID, visit)
}

func (b packageBackend) OpenPath(path string) (Device, error) {
	device, err := sshid.OpenPath(path)
	if err != nil {
		return nil, err
	}
	return &wrappedDevice{device: device}, nil
}

func (d *wrappedDevice) Write(data []byte) (int, error) {
	return d.device.Write(data)
}

func (d *wrappedDevice) Read(data []byte) (int, error) {
	return d.device.Read(data)
}

func (d *wrappedDevice) Close() error {
	return d.device.Close()
}

func openWithBackend(b backend, vendorID, productID uint16) (Device, error) {
	var path string
	stop := errors.New("device selected")
	err := b.Enumerate(vendorID, productID, func(info *sshid.DeviceInfo) error {
		if info == nil || info.Path == "" {
			return nil
		}
		path = info.Path
		return stop
	})
	if err != nil && !errors.Is(err, stop) {
		return nil, err
	}
	if path == "" {
		return nil, ErrDeviceNotFound
	}
	return b.OpenPath(path)
}
