package vnet

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

var ErrNotImplemented = errors.New("not implemented")

const (
	JAIL_CREATE = 0x01
	JAIL_ATTACH = 0x04
)

// Set sets the host or current jail to the jail represented
// by VjHandle.
func Set(vj VjHandle) error {
	_, _, errno := unix.Syscall(unix.SYS_JAIL_ATTACH, uintptr(vj), 0, 0)
	if errno != 0 {
		return fmt.Errorf("jail_attach failed: %s", errno.Error())
	}

	return nil
}

// New creates a new vnet jail, sets it as current and returns
// a handle to it.
func New() (VjHandle, error) {
	iov, err := init_vnet()
	if err != nil {
		return -1, fmt.Errorf("init_vnet failed: %s", err)	
	}

	_, _, errno := unix.Syscall(
		unix.SYS_JAIL_SET, 
		uintptr(unsafe.Pointer(&iov[0])), 
		uintptr(len(iov)),
		uintptr(JAIL_CREATE|JAIL_ATTACH),
	)
	if errno != 0 {
		return -1, fmt.Errorf("jail_set failed: %s", errno.Error())
	} 
	
	return Get()
}

// init_vnet returns []unix.Iovec{} for vnet jail.
func init_vnet() ([]unix.Iovec, error) {
	params := []struct {
		key string
		value interface{}
	}{
		{"path", "/"},
		{"vnet", 1},
		{"persist", nil},
	}

	iovs := []unix.Iovec{}
	for _, param := range params {
		// Add key
		iovs = append(iovs, unix.Iovec{
			Base: (*byte)(unsafe.Pointer(&[]byte(param.key)[0])),
			Len:  uint64(len(param.key) + 1),
		})

		// Add value
		switch v := param.value.(type) {
		case string:
			iovs = append(iovs, unix.Iovec{
				Base: (*byte)(unsafe.Pointer(&[]byte(v)[0])),
				Len:  uint64(len(v) + 1),
			})
		case int:
			iovs = append(iovs, unix.Iovec{
				Base: (*byte)(unsafe.Pointer(&v)),
				Len:  uint64(unsafe.Sizeof(int32(v))),
			})
		case nil:
			iovs = append(iovs, unix.Iovec{
				Base: nil,
				Len: 0,
			})
		default:
			return nil, fmt.Errorf("Unspported vakue type: {%s, %v}\n", param.key, param.value)
		}
	}

	return iovs, nil
}

func NewNamed(name string) (VjHandle, error) {
	return -1, ErrNotImplemented
}

func DeleteNamed(name string) error {
	return ErrNotImplemented
}

func Get() (VjHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromPath(path string) (VjHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromName(name string) (VjHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromPid(pid int) (VjHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromThread(pid int, tid int) (VjHandle, error) {
	return -1, ErrNotImplemented
}

