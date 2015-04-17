package server

import (
//	"eticket-billing-server/operations"
//	"github.com/golang/glog"
)

func Transaction(req *Request) *Request {
	return req.Perform(func(req *Request) string {
/*		transaction := operations.NewTransactionWithoutCheck(req.XmlBody)
		if _, err := transaction.Save(); err != nil {
			response := transaction.ErrorXmlResponse(err)
			glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
			return response
		} else {
			response := transaction.XmlResponse()
			glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
			return response
		}*/
		return ""
	})
}
