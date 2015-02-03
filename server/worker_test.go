package server

import (
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func TestWorker(t *testing.T) { TestingT(t) }

type WorkerSuite struct{}

var _ = Suite(&WorkerSuite{})

func (s *WorkerSuite) TestnewWorker(c *C) {
	worker := newWorker("1", "/tmp")

	c.Assert(reflect.TypeOf(worker).String(), Equals, "*server.Worker")
	c.Assert(worker.merchant, Equals, "1")
}
