package server

type MiddlewareChain func(*Context) *Context

func NewChain(constructors ...func(func(*Context) *Context) func(*Context) *Context) MiddlewareChain {
	var lastFunction MiddlewareChain
	for i := len(constructors) - 1; i >= 0; i-- {
		if i == len(constructors)-1 {
			lastFunction = constructors[i](func(r *Context) *Context {
				return r
			})
		} else {
			lastFunction = constructors[i](lastFunction)
		}
	}
	return lastFunction
}
