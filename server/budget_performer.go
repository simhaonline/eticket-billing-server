package server

import (
	"github.com/golang/glog"
)

type BudgetPerformer struct {
	Request *Request
	Db      *DbConnection
}

func NewBudgetPerformer(request *Request, connection *DbConnection) performerType {
	b := BudgetPerformer{Request: request, Db: connection}
	return performerType(&b)
}

func (p *BudgetPerformer) Serve() *Request {
	return p.Request.Perform(func(req *Request) string {
		budget := NewBudget(req.XmlBody, p.Request.Merchant, p.Db)
		budget.Calculate()
		response := budget.XmlResponse()
		glog.Infof("Worker[%v] answering with %v", p.Request.Merchant, response)
		return response
	})
}
