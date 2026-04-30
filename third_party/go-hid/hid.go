package hid

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	VendorIDAny  uint16 = 0
	ProductIDAny uint16 = 0
)

type Device struct {
	file *os.File
}

type DeviceInfo struct {
	Path         string
	VendorID     uint16
	ProductID    uint16
	SerialNbr    string
	ReleaseNbr   uint16
	MfrStr       string
	ProductStr   string
	UsagePage    uint16
	Usage        uint16
	InterfaceNbr int
}

type EnumFunc func(info *DeviceInfo) error

func Init() error {
	return nil
}

func Exit() error {
	return nil
}

func Enumerate(vid, pid uint16, enumFn EnumFunc) error {
	if enumFn == nil {
		return nil
	}
	entries, err := filepath.Glob("/sys/class/hidraw/hidraw*")
	if err != nil {
		return err
	}
	sort.Strings(entries)
	for _, entry := range entries {
		info, err := readDeviceInfo(entry)
		if err != nil {
			continue
		}
		if vid != VendorIDAny && info.VendorID != vid {
			continue
		}
		if pid != ProductIDAny && info.ProductID != pid {
			continue
		}
		if err := enumFn(info); err != nil {
			return err
		}
	}
	return nil
}

func Open(vid, pid uint16, serial string) (*Device, error) {
	var path string
	err := Enumerate(vid, pid, func(info *DeviceInfo) error {
		if info == nil {
			return nil
		}
		if serial != "" && info.SerialNbr != serial {
			return nil
		}
		path = info.Path
		return io.EOF
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if path == "" {
		return nil, os.ErrNotExist
	}
	return OpenPath(path)
}

func OpenFirst(vid, pid uint16) (*Device, error) {
	return Open(vid, pid, "")
}

func OpenPath(path string) (*Device, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &Device{file: file}, nil
}

func (d *Device) Close() error {
	if d == nil || d.file == nil {
		return nil
	}
	return d.file.Close()
}

func (d *Device) Read(p []byte) (int, error) {
	if d == nil || d.file == nil {
		return 0, os.ErrClosed
	}
	return d.file.Read(p)
}

func (d *Device) Write(p []byte) (int, error) {
	if d == nil || d.file == nil {
		return 0, os.ErrClosed
	}
	return d.file.Write(p)
}

func readDeviceInfo(sysPath string) (*DeviceInfo, error) {
	fields, err := readUEvent(filepath.Join(sysPath, "device", "uevent"))
	if err != nil {
		return nil, err
	}
	id := fields["HID_ID"]
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return nil, errors.New("invalid HID_ID")
	}
	vid, err := parseHex16(parts[1])
	if err != nil {
		return nil, err
	}
	pid, err := parseHex16(parts[2])
	if err != nil {
		return nil, err
	}
	return &DeviceInfo{
		Path:         filepath.Join("/dev", filepath.Base(sysPath)),
		VendorID:     vid,
		ProductID:    pid,
		MfrStr:       fields["HID_VENDOR"],
		ProductStr:   fields["HID_NAME"],
		SerialNbr:    fields["HID_UNIQ"],
		InterfaceNbr: -1,
	}, nil
}

func readUEvent(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fields := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		fields[key] = value
	}
	return fields, nil
}

func parseHex16(value string) (uint16, error) {
	value = strings.TrimSpace(value)
	if len(value) > 4 {
		value = value[len(value)-4:]
	}
	parsed, err := strconv.ParseUint(value, 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(parsed), nil
}
