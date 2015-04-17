package server

import (
	"github.com/golang/glog"
)

type TransactionPerformer struct {
	Request *Request
	Db      *DbConnection
}

func NewTransactionPerformer(request *Request, connection *DbConnection) performerType {
	t := TransactionPerformer{Request: request, Db: connection}
	return performerType(&t)
}

func (p *TransactionPerformer) Serve() *Request {
	return p.Request.Perform(func(req *Request) string {
		transaction := NewTransactionWithoutCheck(req.XmlBody, p.Db)
		transaction.Db.Db.Create(transaction)
		// TODO hold error

		response := transaction.XmlResponse()
		glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
		return response

		/*		transaction := operations.NewTransactionWithoutCheck(req.XmlBody)
				if _, err := transaction.Save(); err != nil {
					response := transaction.ErrorXmlResponse(err)
					glog.Infof("Worker[%v] answering with %v", req.Merchant, response)
					return response
				} else {


					return response
				}*/
	})
}
