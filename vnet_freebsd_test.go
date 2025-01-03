package vnet

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
)

func TestNewGet(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	newns, err := New()
	if err != nil {
		t.Fatal(err)
	}

	ok, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	if newns != ok {
		t.Fatal(fmt.Errorf("newns id = %d, But Get() return %d", newns, ok))
	}
}

func TestSetGet(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	jid, err := strconv.Atoi(os.Getenv("VNET_SET_TEST"))
	err = Set(VjHandle(jid))
	if err != nil {
		t.Fatal(err)
	}

	ok, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	if VjHandle(jid) != ok {
		t.Fatal(fmt.Errorf("jid = %d, But Get() return %d", jid, ok))
	}
}

