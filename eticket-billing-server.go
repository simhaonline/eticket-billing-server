package main

import (
	"eticket-billing-server/config"
	"eticket-billing-server/operations"
	"eticket-billing-server/performers"
	"eticket-billing-server/server"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

	// TODO server should handle it by itself
	pid := syscall.Getpid()
	glog.Info(pid)
	spid := fmt.Sprintf("%v", pid)
	err := ioutil.WriteFile(pidfile, []byte(spid), 0644)
	if err != nil {
		glog.Fatalf("Could not open pidfile. %v", err)
		panic(err)
	}

	mapping := make(performers.PerformerFnMapping)
	mapping["budget"] = performers.Budget
	mapping["transaction"] = performers.Transaction
	performers.SetupMapping(mapping)

	config := config.NewConfig(environment, configFile)
	operations.SetupConnections(config)

	chain := server.NewChain(middleware.NewPingMiddleware, middleware.NewLogMiddleware, middleware.NewServeMiddleware)
	// TODO we can pass mapping to server
	server := server.NewServer(config, chain)
	glog.Infof("New Server is starting with configuration %+v", config)
	glog.Flush()

	go server.Serve()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	signal := <-signalChan

	glog.V(2).Infof("Received %v signal. Stopping...", signal)
	glog.Flush()

	waitForStop := make(chan bool, 1)
	server.Stop(waitForStop)

	<-waitForStop
	glog.Info("EXIT")
}
