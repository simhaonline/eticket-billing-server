package server

import (
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func TestWorker(t *testing.T) { TestingT(t) }

type WorkerSuite struct {
	chain MiddlewareChain
}

var _ = Suite(&WorkerSuite{})

func (s *WorkerSuite) SetUpSuite(c *C) {
	s.chain = NewChain(NewServeMiddleware)
}

func (s *WorkerSuite) TestnewWorker(c *C) {
	worker := newWorker("1", s.chain, "/tmp")

	c.Assert(reflect.TypeOf(worker).String(), Equals, "*server.Worker")
	c.Assert(worker.merchant, Equals, "1")
}
