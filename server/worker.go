package server

import (
    "fmt"
    "os"
    "strconv"
    glog "github.com/golang/glog"
    "eticket-billing/operations"
)

type Worker struct {
    merchant string
    inputChan chan *Request
    quitChan chan bool
    requestsLog *os.File
}

func newWorker(merchant string, filePrefix string) *Worker {
    m, _ := strconv.Atoi(merchant)
    fileName := fmt.Sprintf("%v/worker_%v.log", filePrefix, m)

    f, ok := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if ok != nil {
        glog.Fatal(ok)
        panic(ok)
    }

    return &Worker{merchant, make(chan *Request), make(chan bool), f}
}

func (w Worker) logRequest(req string) {
    _, err := w.requestsLog.WriteString(req + "\n")
    if err != nil {
        glog.Fatal(err)
        panic(err)
    }
}

func (w Worker) Serve() {
    glog.Info("New Worker[%v] is spawned", w.merchant)

    var req *Request
    for {
        select {
        case req = <- w.inputChan:
            w.logRequest(req.XmlBody)
            glog.Infof("Worker[%v] received income request %v", w.merchant, req.XmlBody)

            switch req.OperationType {
            case "budget":
                req.Performer(func(req *Request) string {
                    budget := operations.Budget{Merchant: w.merchant}
                    budget.Calculate()
                    response := budget.XmlResponse()
                    glog.Infof("Worker[%v] answering with %v", w.merchant, response)
                    return response
                })
            case "transaction":
                req.Performer(func(req *Request) string {
                    transaction := operations.NewTransaction(req.XmlBody)
                    if _, err := transaction.Save(); err != nil {
                        response := transaction.ErrorXmlResponse(err)
                        glog.Infof("Worker[%v] answering with %v", w.merchant, response)
                        return response
                    } else {
                        response := transaction.XmlResponse()
                        glog.Infof("Worker[%v] answering with %v", w.merchant, response)
                        return response
                    }
                })
            default:
                glog.Errorf("Worker[%v] received unexpected request %v", w.merchant, req.XmlBody)
                req.Conn.Write([]byte("I have no idea what to do\n"))
                req.Conn.Close()
            }
        case <- w.quitChan:
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
