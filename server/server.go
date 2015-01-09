package server

import (
    "net"
    "os"
    "fmt"
    "io"
    "strings"
    "bytes"
    glog "github.com/golang/glog"
)

var (
    listOfWorkers = make([]*Worker, 0)
)

type Server struct {
    requestLog *os.File
}

func NewServer(inputRequestLogPath string) *Server {
    f, ok := os.OpenFile(inputRequestLogPath + "/requests.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if ok != nil { panic(ok) }

    s := Server{f}
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

    l, err := net.Listen("tcp", ":2000")
    if err != nil {
        glog.Fatal(err)
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            glog.Fatal(err)
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
