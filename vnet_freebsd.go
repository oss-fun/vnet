package vnet

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	return fmt.Sprintf("/var/run/netns/netns%d", int(vj))
}

// Set sets the host or current jail to the jail represented
// by VjHandle.
func Set(vj VjHandle) error {
	_, _, errno := unix.Syscall(unix.SYS_JAIL_ATTACH, uintptr(vj), 0, 0)
	if errno != 0 {
		return fmt.Errorf("jail_attach failed: %s", errno.Error())
	}

	f, err := os.OpenFile(vnetPath(vj), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("init_vnet failed: %s", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%d\n", os.Getpid())); err != nil {
		return fmt.Errorf("WriteString failed: %s", err)
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

	f, err := os.OpenFile(vnetPath(VjHandle(jid)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return VjHandle(jid), fmt.Errorf("OpenFile failed: %s", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%d\n", os.Getpid())); err != nil {
		return VjHandle(jid), fmt.Errorf("WriteString failed: %s", err)
	}

	return VjHandle(jid), nil
}

// init_vnet returns []unix.Iovec{} for vnet jail.
func init_vnet() ([]unix.Iovec, error) {
	params := []struct {
		key   string
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
				Len:  0,
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
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(path)
	if match == "" {
		return -1, fmt.Errorf("no trailing number found")
	}

	jid, err := strconv.Atoi(match)
	if err != nil {
		return -1, err
	}
	return VjHandle(jid), nil
}

func GetFromName(name string) (VjHandle, error) {
	return -1, ErrNotImplemented
}

// GetFromPid gets a handle to the vnet jail of a given pid.
func GetFromPid(pid int) (VjHandle, error) {
	dir := "/var/run/netns"
	var result string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if num, err := strconv.Atoi(line); err == nil && num == pid {
				result = path
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return -1, err
	}

	return GetFromPath(result)
}
