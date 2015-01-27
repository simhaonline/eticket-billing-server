package server

func NewLogMiddleware(f func(*Request)  *Request) func(*Request) *Request {
    return func(req *Request) *Request {
        return f(req)
    }
}
