package server

func NewLogMiddleware(f func(*Context) *Context) func(*Context) *Context {
	return func(context *Context) *Context {
		return f(context)
	}
}
