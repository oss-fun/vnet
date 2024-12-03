//go:build !freebsd
// +build !freebsd

package vnet

import "errors"

var ErrNotImplemented = errors.New("not implemented")

func SetVnet(vj VjHandle, vjtype int) error {
	return ErrNotImplemented
}

func Set(vj VjHandle) error {
	return ErrNotImplemented
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

