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

/*
func TestClose(t *testing.T) {
	var jid VjHandle
	jid := 0
	err := jid.Close()
	if err != nil {
		t.Fatal(err)
	}
}
*/

/*
func TestString(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var jid VjHandle
	jid = 0
	fmt.Println(jid.String())
}
*/
