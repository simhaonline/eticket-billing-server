package server

type Context struct {
	Request           *Request
	Db                *DbConnection
	PerformersMapping *PerformerFnMapping
}
