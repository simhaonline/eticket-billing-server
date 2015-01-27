package server

import(
    "eticket-billing-server/operations"
    "github.com/golang/glog"
)

func NewServeMiddleware(f func(*Request) *Request) func(*Request) *Request {
    return func(req *Request) *Request {

        switch req.OperationType {
        case "budget":
            req.Performer(func(req *Request) string {
                budget := operations.Budget{Merchant: req.Merchant}
                budget.Calculate()
                response := budget.XmlResponse()
                glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
                return response
            })
        case "transaction":
            req.Performer(func(req *Request) string {
                transaction := operations.NewTransaction(req.XmlBody)
                if _, err := transaction.Save(); err != nil {
                    response := transaction.ErrorXmlResponse(err)
                    glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
                    return response
                } else {
                    response := transaction.XmlResponse()
                    glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
                    return response
                }
            })
        case "transaction-without-check":
            req.Performer(func(req *Request) string {
                transaction := operations.NewTransactionWithoutCheck(req.XmlBody)
                if _, err := transaction.Save(); err != nil {
                    response := transaction.ErrorXmlResponse(err)
                    glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
                    return response
                } else {
                    response := transaction.XmlResponse()
                    glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
                    return response
                    }
            })
        default:
            glog.Errorf("Worker[%v] received unexpected request %v", req.Merchant, req.XmlBody)
            req.Conn.Write([]byte("I have no idea what to do\n"))
            req.Conn.Close()
        }

        return f(req)
    }
}
