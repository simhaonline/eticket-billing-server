package server

import (
	"eticket-billing-server/config"
	glog "github.com/golang/glog"
	"bytes"
	"io"
	"net"
	"strings"
	"time"
	"flag"
	"io/ioutil"
	"os"
	"syscall"
	"fmt"
)

var (
	listOfWorkers = make([]*Worker, 0)
)

type Server struct {
	stopChan    chan bool
	requestLog  *os.File
	config      *config.Config
	middlewares MiddlewareChain
	workersPool WorkersPool
	performersMapping PerformerFnMapping

	pidFile string
}

func NewServer(config *config.Config, middlewares MiddlewareChain, mapping PerformerFnMapping) *Server {
	f, ok := os.OpenFile(config.RequestLogDir+"/requests.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if ok != nil {
		panic(ok)
	}

	s := Server{stopChan: make(chan bool), requestLog: f, config: config, middlewares: middlewares, performersMapping: mapping}
	s.prepareServer()
	return &s
}

func (s Server) logRequest(req string) {
	_, err := s.requestLog.WriteString(req + "\n")
	if err != nil {
		glog.Fatal(err)
		panic(err)
	}
}

func (s *Server) Serve() {
	s.writePid()
	glog.Infof("New Server is started with configuration %+v", s.config)
	glog.Flush()

	laddr, err := net.ResolveTCPAddr("tcp", ":2000")
	if nil != err {
		glog.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		glog.Fatal(err)
	}
	defer l.Close()

	s.workersPool = NewWorkersPool(s.config, s.middlewares, s.performersMapping)

	for {
		select {
		case <-s.stopChan:
			glog.Info("Breaking listening loop. Closing Connection...")
			glog.Flush()
			return
		default:

		}
		l.SetDeadline(time.Now().Add(1e9))
		conn, err := l.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			glog.Info(err)
		}

		conn.SetDeadline(time.Now().Add(1e9))
		buf := make([]byte, 1024)
		if _, err := conn.Read(buf); nil != err && err != io.EOF {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			glog.Fatal(err)
			return
		}

		buf = bytes.Trim(buf, "\x00")

		if len(buf) == 0 {
			continue
		}

		input := string(buf)
		input = strings.TrimSpace(input)

		request := NewRequest(input)
		request.Conn = conn

		glog.Info(input)

		s.logRequest(input)

		worker := s.workersPool.GetWorkerForMerchant(request.Merchant)
		worker.inputChan <- request
	}
}

func (s *Server) Stop(stChan chan bool) {
	glog.Info("Attempting to stop everything")
	s.stopChan <- true
	close(s.stopChan)
	s.requestLog.Close()

	glog.V(2).Info("Closed servers files and chans")
	glog.Flush()

	s.workersPool.StopAll()

	glog.Info("Server is stopped")
	glog.Flush()

	stChan <- true
}

func (s *Server) prepareServer() {
	var environment string
	var configFile string

	flag.StringVar(&environment, "environment", "", "Setup environment: production, development")
	flag.StringVar(&configFile, "config-file", "", "Configuration file")
	flag.StringVar(&s.pidFile, "pidfile", "", "PID file")
	flag.Parse()

	s.config = config.NewConfig(environment, configFile)
}

func (s *Server) writePid() {
	pid := syscall.Getpid()
	glog.Info(pid)
	spid := fmt.Sprintf("%v", pid)
	err := ioutil.WriteFile(s.pidFile, []byte(spid), 0644)
	if err != nil {
		glog.Fatalf("Could not open pidfile. %v", err)
		panic(err)
	}
}
