package vnet

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

var ErrNotImplemented = errors.New("not implemented")

const (
	JAIL_CREATE = 0x01
	JAIL_ATTACH = 0x04
)


// vnetPath return /var/run/netns/netns<jid>.
func vnetPath(vj VjHandle) string {
	return fmt.Sprintf("/var/run/netns/netns%d", vj)
}

// status return /proc/<pid>/status.
func status(pid int) string {
	return fmt.Sprintf("/proc/%d/status", pid)
}

// Set sets the host or current jail to the jail represented
// by VjHandle.
func Set(vj VjHandle) error {
	_, _, errno := unix.Syscall(unix.SYS_JAIL_ATTACH, uintptr(vj), 0, 0)
	if errno != 0 {
		return fmt.Errorf("jail_attach failed: %s", errno.Error())
	}

	if err := os.Symlink(status(os.Getpid()), vnetPath(vj)); err != nil {
		return fmt.Errorf("Symlink failed: %s", err)
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

	jid, _, errno := unix.Syscall(
		unix.SYS_JAIL_SET, 
		uintptr(unsafe.Pointer(&iov[0])), 
		uintptr(len(iov)),
		uintptr(JAIL_CREATE|JAIL_ATTACH),
	)
	if errno != 0 {
		return -1, fmt.Errorf("jail_set failed: %s", errno.Error())
	}

	if err := os.Symlink(status(os.Getpid()), vnetPath(VjHandle(jid))); err != nil {
		return VjHandle(jid), fmt.Errorf("Symlink failed: %s", err)
	}

	return VjHandle(jid), nil
}

// init_vnet returns []unix.Iovec{} for vnet jail.
func init_vnet() ([]unix.Iovec, error) {
	params := []struct {
		key string
		value interface{}
	}{
		{"path", "/"},
		{"vnet", 1},
		{"children.max", 99},
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
			return nil, fmt.Errorf("Unspported value type: {%s, %v}", param.key, param.value)
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

// Get gets a handle to the current vnet jail.
func Get() (VjHandle, error) {
	return GetFromPid(os.Getpid())
}

// GetFromPath gets a jail ID to a network namespace
// identidied by the path
func GetFromPath(path string) (VjHandle, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf("Open failed: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("scanner failed: %s", err)
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) == 0 {
		return -1, fmt.Errorf("no fields found in the last line")
	}
	if fields[len(fields)-1] == "-" {
		return 0, fmt.Errorf("The process specified by the path is running on the host.")
	}
	id, err := strconv.Atoi(fields[len(fields)-1])
	if err != nil {
		return -1, fmt.Errorf("Atoi failed: %s", err)
	}
	return VjHandle(id), nil
}

func GetFromName(name string) (VjHandle, error) {
	return -1, ErrNotImplemented
}

// GetFromPid gets a handle to the vnet jail of a given pid.
func GetFromPid(pid int) (VjHandle, error) {
	return GetFromPath(fmt.Sprintf("/proc/%d/status", pid))
}

