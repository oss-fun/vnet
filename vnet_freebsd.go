package vnet

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

var ErrNotImplemented = errors.New("not implemented")

const (
	JAIL_CREATE = 0x01
	JAIL_ATTACH = 0x04
)

func Set(vj VjHandle) error {
	_, _, errno := unix.Syscall(uintptr(unix.SYS_JAIL_ATTACH), uint(jid), 0, 0)
	if errno != 0 {
		return fmt.Errorf("Jail_attach: %s", errno.Error())
	}

	return nil
}

func New() (VjHandle, error) {
	return -1, ErrNotImplemented
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

