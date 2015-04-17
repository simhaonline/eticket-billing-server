package server

import (
	"eticket-billing-server/config"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func TestWorker(t *testing.T) { TestingT(t) }

type WorkerSuite struct {
	chain             MiddlewareChain
	config            *config.Config
	performersMapping *PerformerFnMapping
}

var _ = Suite(&WorkerSuite{})

func (s *WorkerSuite) SetUpSuite(c *C) {
	s.chain = NewChain(NewServeMiddleware)

	mapping := make(PerformerFnMapping)
	mapping["budget"] = NewBudgetPerformer

	s.config = config.NewConfig("test", "../config.gcfg")
	s.performersMapping = &mapping
}

func (s *WorkerSuite) TestnewWorker(c *C) {
	// worker := newWorker("1", s.chain, "/tmp")
	worker := newWorker("1", s.chain, s.config, s.performersMapping)

	c.Assert(reflect.TypeOf(worker).String(), Equals, "*server.Worker")
	c.Assert(worker.merchant, Equals, "1")
}
