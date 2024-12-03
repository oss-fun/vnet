# vnet - network namespace conpatible vnet in go #

The vnet package provides an ultra-simple interface for handling vnet in go. 
Useing vnet jail requires elevated privileges, so in most cases this code needs to be run as root.

## Local Build and Test ##

You can use go get command:

    go get github.com/oss-fun/vnet

Testing (requires root):

    sudo -E go test github.com/oss-fun/vnet

## Example ##

```go
package main

import (
    "fmt"
    "net"
    "runtime"

    "github.com/vishvananda/netns"
)

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := vnet.Get()
    defer origns.Close()

    // Create a new network namespace
    newns, _ := vnet.New()
    defer newns.Close()

    // Do something with the network namespace
    ifaces, _ := net.Interfaces()
    fmt.Printf("Interfaces: %v\n", ifaces)

    // Switch back to the original namespace
    netns.Set(origns)
}

```
