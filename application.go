package main

import (
    glog "github.com/golang/glog"
    "flag"
    "eticket-billing/server"
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "time"
)

func main() {
    flag.Parse()
    defer glog.Flush()

    server := server.NewServer("/tmp/")
    go server.Serve()

    signalChan := make(chan os.Signal)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

    fmt.Println(<-signalChan)
    server.Stop()

    glog.Info("EXIT")
    glog.Flush()
    time.Sleep(100 * time.Millisecond)
    os.Exit(0)
}
