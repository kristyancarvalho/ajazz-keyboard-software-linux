package hid

import (
	"errors"
	"testing"

	sshid "github.com/sstallion/go-hid"
	"github.com/stretchr/testify/require"
)

type fakeBackend struct {
	infos []*sshid.DeviceInfo
	open  Device
	err   error
}

func (b fakeBackend) Enumerate(vendorID, productID uint16, visit func(info *sshid.DeviceInfo) error) error {
	if b.err != nil {
		return b.err
	}
	for _, info := range b.infos {
		if err := visit(info); err != nil {
			return err
		}
	}
	return nil
}

func (b fakeBackend) OpenPath(path string) (Device, error) {
	if b.open == nil {
		return nil, errors.New("no device")
	}
	return b.open, nil
}

func TestOpenReturnsDeviceNotFound(t *testing.T) {
	t.Parallel()
	_, err := openWithBackend(fakeBackend{}, 1, 2)
	require.ErrorIs(t, err, ErrDeviceNotFound)
}

func TestFakeDeviceRecordsWrites(t *testing.T) {
	t.Parallel()
	device := &FakeDevice{}
	_, err := device.Write([]byte{1, 2, 3})
	require.NoError(t, err)
	require.Len(t, device.Writes, 1)
	require.Equal(t, []byte{1, 2, 3}, device.Writes[0])
}
