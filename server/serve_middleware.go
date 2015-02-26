package server

import (
	"eticket-billing-server/performers"
)

func NewServeMiddleware(f func(*request.Request) *request.Request) func(*request.Request) *request.Request {
	return func(req *request.Request) *request.Request {
		handler := performers.GetMapping(req.OperationType)
		handler(req)
		return f(req)
	}
}
