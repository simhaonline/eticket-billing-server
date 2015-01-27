package server

import(
    "eticket-billing-server/request"
)

type MiddlewareChain func(*request.Request) *request.Request

func NewChain(constructors ...func(func(*request.Request) *request.Request) func(*request.Request) *request.Request) MiddlewareChain {
    var lastFunction MiddlewareChain
    for i := len(constructors)-1; i >=0 ; i-- {
        if i == len(constructors)-1 {
            lastFunction = constructors[i](func(r *request.Request) *request.Request {
                return r
            })
        } else {
            lastFunction = constructors[i](lastFunction)
        }
    }
    return lastFunction
}
