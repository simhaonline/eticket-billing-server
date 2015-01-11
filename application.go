// +build main

package main

import (
    "github.com/golang/glog"
    "flag"
    "eticket-billing/server"
    "eticket-billing/config"
    "os"
    "os/signal"
    "syscall"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    defer glog.Flush()

    var environment string
    var configFile string

    flag.StringVar(&environment, "environment", "", "Setup environment: production, development")
    flag.StringVar(&configFile, "config-file", "", "Configuration file")
    flag.Parse()

    config.ParseConfig(environment, configFile)

    waitForStop := make(chan bool, 1)

    config := config.GetConfig()

    server := server.NewServer(config)
    glog.Infof("New Server is starting with configuration %+v", config)
    glog.Flush()
    go server.Serve()

    signalChan := make(chan os.Signal)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

    signal := <-signalChan

    glog.V(2).Infof("Received %v signal. Stopping...", signal)
    glog.Flush()

    server.Stop(waitForStop)

    <-waitForStop
    glog.Info("EXIT")
}
