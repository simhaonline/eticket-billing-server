package server

type MiddlewareChain func(*Request) *Request

func NewChain(constructors ...func(func(*Request) *Request) func(*Request) *Request) MiddlewareChain {
    var lastFunction MiddlewareChain
    for i := len(constructors)-1; i >=0 ; i-- {
        if i == len(constructors)-1 {
            lastFunction = constructors[i](func(r *Request) *Request {
                return r
            })
        } else {
            lastFunction = constructors[i](lastFunction)
        }
    }
    return lastFunction
}
