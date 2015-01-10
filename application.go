package main

import (
    glog "github.com/golang/glog"
    "flag"
    "eticket-billing/server"
    "os"
    "os/signal"
    "syscall"
    "path/filepath"
)

func main() {
    defer glog.Flush()

    config := server.NewConfig()

    waitForStop := make(chan bool, 1)

    flag.StringVar(&config.Environment, "environment", "development", "Setup environemnt: production, development")
    flag.StringVar(&config.RequestLogDir, "request-log-dir", "", "Place where to store requests log files")
    flag.Parse()

    if config.RequestLogDir == "" {
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            glog.Fatal(err)
            os.Exit(1)
        }

        config.RequestLogDir = dir
        glog.Infof("Use directory `%v' as root for storing log files", dir)
    }

    config.DataBaseName = "eticket_billing_" + config.Environment

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
