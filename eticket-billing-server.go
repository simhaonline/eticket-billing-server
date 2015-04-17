package main

import (
	"eticket-billing-server/operations"
	"eticket-billing-server/server"
	"fmt"
	"github.com/golang/glog"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer glog.Flush()

	mapping := make(server.PerformerFnMapping)
	mapping["budget"] = server.NewBudgetPerformer
	mapping["transaction"] = server.NewTransactionPerformer

	chain := server.NewChain(server.NewPingMiddleware, server.NewLogMiddleware, server.NewServeMiddleware)

	server := server.NewServer(chain, mapping)

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
