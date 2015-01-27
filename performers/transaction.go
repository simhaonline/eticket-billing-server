package performers

import(
    "eticket-billing-server/operations"
    "eticket-billing-server/request"
    "github.com/golang/glog"
)

func Transaction(req *request.Request) *request.Request {
    return req.Perform(func(req *request.Request) string {
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
}
