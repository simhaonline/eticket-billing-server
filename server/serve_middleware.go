package server

func NewServeMiddleware(f func(*Context) *Context) func(*Context) *Context {
	return func(context *Context) *Context {
		defer context.Db.Db.Close()

		performerConstructor := (*context.PerformersMapping)[context.Request.OperationType]
		performer := performerConstructor(context.Request, context.Db)
		_ = (performerType(performer)).Serve()
		return f(context)
	}
}
