package vnet

import (
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}
}
