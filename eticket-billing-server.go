package main

import (
    "github.com/golang/glog"
    "flag"
    "eticket-billing-server/server"
    "eticket-billing-server/config"
    "os"
    "os/signal"
    "syscall"
    "runtime"
    "io/ioutil"
    "fmt"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    defer glog.Flush()

    var environment string
    var configFile string
    var pidfile string

    flag.StringVar(&environment, "environment", "", "Setup environment: production, development")
    flag.StringVar(&configFile, "config-file", "", "Configuration file")
    flag.StringVar(&pidfile, "pidfile", "", "PID file")
    flag.Parse()

    pid := syscall.Getpid()
    glog.Info(pid)
    spid := fmt.Sprintf("%v", pid)
    err := ioutil.WriteFile(pidfile, []byte(spid), 0644)
    if err != nil {
        glog.Fatalf("Could not open pidfile. %v", err)
        panic(err)
    }

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
