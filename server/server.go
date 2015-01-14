package server

import (
    "net"
    "os"
    "io"
    "strings"
    "bytes"
    "time"
    "github.com/golang/glog"
    "eticket-billing/config"
)

var (
    listOfWorkers = make([]*Worker, 0)
)

type Server struct {
    stopChan chan bool
    requestLog *os.File
    config *config.Config
}

func NewServer(config *config.Config) *Server {
    // TODO check connection to DB
    f, ok := os.OpenFile(config.RequestLogDir + "/requests.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if ok != nil { panic(ok) }

    s := Server{stopChan: make(chan bool), requestLog: f, config: config}
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
    glog.Info("Ready")

    laddr, err := net.ResolveTCPAddr("tcp", ":2000")
    if nil != err {
        glog.Fatal(err)
    }

    l, err := net.ListenTCP("tcp", laddr)
    if err != nil {
        glog.Fatal(err)
    }
    defer l.Close()

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
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout()  {
                continue
            }
            glog.Fatal(err)
            return
        }

        buf = bytes.Trim(buf, "\x00")

        if len(buf) == 0 { continue }

        input := string(buf)
        input = strings.TrimSpace(input)

        request := NewRequest(input)
        request.Conn = conn

        glog.Info(input)

        s.logRequest(input)

        pool := GetWorkersPool()

        worker := pool.GetWorkerForMerchant(request.Merchant)
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

    pool := GetWorkersPool()
    pool.StopAll()

    glog.Info("Server is stopped")
    glog.Flush()

    stChan <- true
}
