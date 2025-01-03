package vnet

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// NsHandle is a handle to a vnet jail. It can be cast directly
// to an int and used as a jail ID.
type VjHandle int

// Equal determines if two vnet handles refer to the same vnet jail.
// This is done by comparing the device and inode that the jail point to.
func (vj VjHandle) Equal(other VjHandle) bool {
	if vj == other {
		return true
	}

	var s1, s2 unix.Stat_t
	f1, err := os.Open(vnetPath(vj))
	if err != nil {
		return false
	}
	defer f1.Close()
	f2, err := os.Open(vnetPath(other))
	if err != nil {
		return false
	}
	defer f2.Close()
	if err = unix.Fstat(int(f1.Fd()), &s1); err != nil {
		return false
	}
	if err = unix.Fstat(int(f2.Fd()), &s2); err != nil {
		return false
	}

	return (s1.Dev == s2.Dev) && (s1.Ino == s2.Ino)
}

// String shows the jail ID and dev and inode.
func (vj VjHandle) String() string {
	if vj == -1 {
		return "vnet(none)"
	}
	var s unix.Stat_t
	f, err := os.Open(vnetPath(vj))
	if err != nil {
		return "vnet(none)"
	}
	defer f.Close()
	if err := unix.Fstat(int(f.Fd()), &s); err != nil {
		return "vnet(unknown)"
	}
	return fmt.Sprintf("vnet(%d: %d, %d)", vj, s.Dev, s.Ino)
}

// UniqueId returns a string which uniquely identifies the jail
// associated with the VjHandle.
func (vj VjHandle) UniqueId() string {
	if vj == -1 {
		return "vnet(none)"
	}
	var s unix.Stat_t
	f, err := os.Open(vnetPath(vj))
	if err != nil {
		return "vnet(none)"
	}
	defer f.Close()
	if err := unix.Fstat(int(f.Fd()), &s); err != nil {
		return "vnet(unknown)"
	}
	return fmt.Sprintf("vnet(%d:%d)", s.Dev, s.Ino)
}

// IsOpen returns true if Close() has not been called.
func (vj VjHandle) IsOpen() bool {
	return vj != -1
}

// Close closes the NsHandle and resets its file descriptor to -1.
// In FreeBSD, Close must be called after the process in jail.
func (vj VjHandle) Close() error {
	_, _, errno := unix.Syscall(unix.SYS_JAIL_REMOVE, uintptr(int(vj)), 0, 0)
	if errno != 0 {
		return fmt.Errorf("jail_remove failed: %s", errno.Error())
	}

	err := os.Remove(vnetPath(vj))
	if err != nil {
		return err
	}
	vj = -1
	return nil
}

// None gets an empty (closed) NsHandle.
func None() VjHandle {
	return VjHandle(-1)
}
