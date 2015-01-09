package main

import (
    glog "github.com/golang/glog"
    "flag"
    "eticket-billing/server"
)

func main() {
    flag.Parse()
    defer glog.Flush()

    server := server.NewServer("/tmp/")
    server.Serve()
}
