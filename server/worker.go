package server

import (
    "log"
    "fmt"
    "os"
    "strconv"
)

type Worker struct {
    merchant int
    inputChan chan *Request
    quitChan chan bool
    logger *log.Logger
}

func newWorker(merchant int, filePrefix string) *Worker {
    fileName := fmt.Sprintf("%v/worker_%v.log", filePrefix, strconv.Itoa(merchant))

    f, ok := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if ok != nil { panic(ok) }

    logger := log.New(f, "", log.Ldate|log.Ltime)

    return &Worker{merchant, make(chan *Request), make(chan bool), logger}
}

func (w Worker) Serve() {
    Info.Printf("New worker for merchant %v is spawned", w.merchant)

    var req *Request
    for {
        select {
        case req = <- w.inputChan:
            Info.Printf("Worker of merchant %v received string '%v'", w.merchant, req.body)
            w.logger.Println(req.body)
            req.connection.Write([]byte("Confirm\n"))
            req.connection.Close()
        case <- w.quitChan:
            Info.Println("Wroker for %v quitting", w.merchant)
            return
        }
    }
}
