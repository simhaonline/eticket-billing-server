package server

import (
    "log"
    "fmt"
    "os"
    "strconv"
    "eticket-billing/operations"
)

type Worker struct {
    merchant string
    inputChan chan *Request
    quitChan chan bool
    logger *log.Logger
}

func newWorker(merchant string, filePrefix string) *Worker {
    m, _ := strconv.Atoi(merchant)
    fileName := fmt.Sprintf("%s/worker_%s.log", filePrefix, m)

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
            switch req.OperationType {
            case "budget":
                budget := operations.Budget{Merchant: w.merchant}
                amount, _ := budget.Calculate()
                amountResponse := strconv.AppendInt(make([]byte, 0), amount, 10)
                amountResponse = append(amountResponse, '\n')
                req.Conn.Write(amountResponse)
            default:
                req.Conn.Write([]byte("I have no idea what to do\n"))
            }

            req.Conn.Close()
        case <- w.quitChan:
            Info.Println("Wroker for %v quitting", w.merchant)
            return
        }
    }
}
