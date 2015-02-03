package middleware

import (
	"eticket-billing-server/request"
)

func NewPingMiddleware(f func(*request.Request) *request.Request) func(*request.Request) *request.Request {
	return func(req *request.Request) *request.Request {
		return f(req)
	}
}
