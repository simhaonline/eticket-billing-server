package server

import (
    "net"
    "os"
    "fmt"
    "io"
    "strings"
    "bytes"
    "time"
    glog "github.com/golang/glog"
)

var (
    listOfWorkers = make([]*Worker, 0)
)

type Server struct {
    stopChan chan bool
    requestLog *os.File
}

func NewServer(inputRequestLogPath string) *Server {
    f, ok := os.OpenFile(inputRequestLogPath + "/requests.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if ok != nil { panic(ok) }

    s := Server{make(chan bool), f}
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

    laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:2000")
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
            glog.Info("Stopping server. Closing Connection...")
            break
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

        buf := make([]byte, 1024)
        _, err = conn.Read(buf)
        if err != nil && err != io.EOF { panic(err) }

        buf = bytes.Trim(buf, "\x00")

        if len(buf) == 0 { continue }

        input := string(buf)
        input = strings.TrimSpace(input)

        fmt.Printf("%q+", input)

        request := NewRequest(input)
        request.Conn = conn

        glog.Info(input)

        s.logRequest(input)

        pool := NewWorkersPool()

        worker := pool.GetWorkerForMerchant(request.Merchant)
        worker.inputChan <- request
    }
}

func (s *Server) Stop() {
    glog.Info("Attempting to stop everything")
    s.stopChan <- true
    close(s.stopChan)
    s.requestLog.Close()

    pool := NewWorkersPool()
    pool.StopAll()

    glog.Info("Server is stopped")
    glog.Flush()
}
