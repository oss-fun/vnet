# vnet - network namespace conpatible vnet in go #

The vnet package provides an ultra-simple interface for handling vnet in go. 
Useing vnet jail requires elevated privileges, so in most cases this code needs to be run as root.

## Local Build and Test ##

You can use go get command:

    go get github.com/oss-fun/vnet

Testing (requires root):

    sudo -E go test github.com/oss-fun/vnet -run Test***

## Example ##

```go
package main

import (
        "fmt"
        "net"

        "github.com/oss-fun/vnet"
)

func main() {
        // Create a new network namespace
        newvj, _ := vnet.New()

        // Do something with the network namespace
        curvj, _ := vnet.Get()
        if newvj.Equal(curvj) {
                fmt.Printf("OK!\n")
        }

        ifaces, _ := net.Interfaces()
        fmt.Printf("Interfaces: %v\n", ifaces)
}

```
