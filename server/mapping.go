package server

type performerType interface {
	Serve() *Request
}

type PerformerFn func(*Request, *DbConnection) performerType
type PerformerFnMapping map[string]PerformerFn

type performer interface {
	Serve(req *Request) *Request
}
