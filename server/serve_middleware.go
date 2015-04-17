package server

func NewServeMiddleware(f func(*Request) *Request) func(*Request) *Request {
	return func(req *Request) *Request {
		handler := GetMapping(req.OperationType)
		handler(req)
		return f(req)
	}
}
