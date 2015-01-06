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
                req.Performer(func(req *Request) {
                    budget := operations.Budget{Merchant: w.merchant}
                    budget.Calculate()
                    answer := []byte(budget.XmlResponse())
                    req.Conn.Write(answer)
                })
            case "transaction":
                req.Performer(func(req *Request) {
                    transaction := operations.NewTransaction(req.XmlBody)
                    if transaction.IsPossible() {
                        transaction.Save()
                        answer := []byte(transaction.XmlResponse())
                        req.Conn.Write(answer)
                    } else {
                        req.Conn.Write([]byte("I have not enough money\n"))
                    }
                })
            default:
                req.Conn.Write([]byte("I have no idea what to do\n"))
                req.Conn.Close()
            }
        case <- w.quitChan:
            Info.Println("Wroker for %v quitting", w.merchant)
            return
        }
    }
}
