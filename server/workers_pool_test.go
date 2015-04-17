package server

import (
	"eticket-billing-server/config"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func TestWorkerPool(t *testing.T) { TestingT(t) }

type WorkersPoolSuite struct {
	config *config.Config
	chain  MiddlewareChain
}

var _ = Suite(&WorkersPoolSuite{})

func (s *WorkersPoolSuite) SetUpSuite(c *C) {
	s.config = &config.Config{RequestLogDir: "/tmp"}
	s.chain = NewChain(NewServeMiddleware)
}

func (s *WorkersPoolSuite) TestNewWorkersPool(c *C) {
	pool := NewWorkersPool(s.config, s.chain)
	c.Assert(reflect.TypeOf(pool).String(), Equals, "server.WorkersPool")
	c.Assert(len(pool.pool), Equals, 0)
}

func (s *WorkersPoolSuite) TestWorkersPoolInstance(c *C) {
	pool := NewWorkersPool(s.config, s.chain)

	worker := pool.GetWorkerForMerchant("10")
	c.Assert(worker.merchant, Equals, "10")
	c.Assert(len(pool.pool), Equals, 1)

	worker = pool.GetWorkerForMerchant("10")

	c.Assert(worker.merchant, Equals, "10")
	c.Assert(len(pool.pool), Equals, 1)

	worker = pool.GetWorkerForMerchant("20")
	c.Assert(worker.merchant, Equals, "20")
	c.Assert(len(pool.pool), Equals, 2)
}

func (s *WorkersPoolSuite) TestTwoPools(c *C) {
	pool1 := NewWorkersPool(s.config, s.chain)
	pool2 := NewWorkersPool(s.config, s.chain)

	_ = pool1.GetWorkerForMerchant("10")
	_ = pool1.GetWorkerForMerchant("20")
	_ = pool2.GetWorkerForMerchant("20")

	c.Assert(len(pool1.pool), Equals, 2)
	c.Assert(len(pool2.pool), Equals, 1)
}
