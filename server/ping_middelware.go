package server

func NewPingMiddleware(f func(*Context) *Context) func(*Context) *Context {
	return func(context *Context) *Context {
		return f(context)
	}
}
