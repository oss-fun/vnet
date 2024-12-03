package vnet

type VjHandle int

func (vj VjHandle) Equal(_ VjHandle) bool {
	return false
}

// String shows the file descriptor number and its dev and inode.
// It is only implemented on Linux, and returns "NS(none)" on other
// platforms.
func (vj VjHandle) String() string {
	return "NS(none)"
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
