package vnet

// NsHandle is a handle to a vnet jail. It can be cast directly
// to an int and used as a jail ID.
type VjHandle int

// Equal determines if two vnet handles refer to the same vnet jail.
func (vj VjHandle) Equal(other VjHandle) bool {
	if vj == other {
		return true
	}
	return false
}

// String shows the jail ID and its dev and inode.
func (vj VjHandle) String() string {
	if vj == -1 {
		return "NS(none)"
	}
	return "NS()"
}

// UniqueId returns a string which uniquely identifies the namespace
// associated with the network handle. It is only implemented on Linux,
// and returns "NS(none)" on other platforms.
func (vj VjHandle) UniqueId() string {
	return "NS(none)"
}

// IsOpen returns true if Close() has not been called. It is only implemented
// on Linux and always returns false on other platforms.
func (vj VjHandle) IsOpen() bool {
	return false
}

// Close closes the NsHandle and resets its file descriptor to -1.
// It is only implemented on Linux.
func (vj *VjHandle) Close() error {
	return nil
}

// None gets an empty (closed) NsHandle.
func None() VjHandle {
	return VjHandle(-1)
}

