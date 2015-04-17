package server

type performer interface {
	Serve(req *Request) *Request
}

type performerType interface {
	Serve() *Request
}
