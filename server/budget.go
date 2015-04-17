package server

import (
	"eticket-billing-server/operations"
	"github.com/golang/glog"
)

func Budget(req *Request) *Request {
	return req.Perform(func(req *Request) string {
		budget := operations.Budget{Merchant: req.Merchant}
		budget.Calculate()
		response := budget.XmlResponse()
		glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
		return response
	})
}
