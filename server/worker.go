package server

import (
	"fmt"
	glog "github.com/golang/glog"
	"os"
	"strconv"
)

type Worker struct {
	merchant    string
	inputChan   chan *request.Request
	quitChan    chan bool
	requestsLog *os.File
	middleware  MiddlewareChain
}

func newWorker(merchant string, middleware MiddlewareChain, filePrefix string) *Worker {
	m, _ := strconv.Atoi(merchant)
	fileName := fmt.Sprintf("%v/worker_%v.log", filePrefix, m)

	f, ok := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if ok != nil {
		glog.Fatal(ok)
		panic(ok)
	}

	return &Worker{merchant, make(chan *request.Request), make(chan bool), f, middleware}
}

func (w Worker) logRequest(req string) {
	_, err := w.requestsLog.WriteString(req + "\n")
	if err != nil {
		glog.Fatal(err)
		panic(err)
	}
}

func (w Worker) Serve() {
	glog.Infof("New Worker[%v] is spawned", w.merchant)

	var req *request.Request
	for {
		select {
		case req = <-w.inputChan:
			w.logRequest(req.XmlBody)
			glog.Infof("Worker[%v] received income request %v", w.merchant, req.XmlBody)

			// TODO rename to Serve
			w.middleware(req)

		case <-w.quitChan:
			glog.Infof("Wroker %v quitting", w.merchant)
			return
		}
	}
}

func (w Worker) Stop() {
	w.quitChan <- true
	w.requestsLog.Close()
	close(w.quitChan)
	close(w.inputChan)
	glog.V(2).Infof("Worker[%v] is stopped", w.merchant)
	glog.Flush()
}
