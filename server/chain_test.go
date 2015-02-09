package server

import (
	"eticket-billing-server/request"
	"eticket-billing-server/middleware"
	"eticket-billing-server/performers"
	. "gopkg.in/check.v1"
	"testing"
)

func TestChain(t *testing.T) { TestingT(t) }

type ChainSuite struct {
	CheckArray []string
}

var _ = Suite(&ChainSuite{CheckArray: []string{}})

func BarePerformere(req *request.Request) *request.Request {
	return req;
}


func (s *ChainSuite) SetUpSuite(c *C) {
	mapping := make(performers.PerformerFnMapping)
	mapping["bare"] = BarePerformere
	performers.SetupMapping(mapping)
}

func (s *ChainSuite) TestNewChain(c *C) {
	request := request.Request{Merchant: "m10", OperationType: "bare"}
	chain := NewChain(middleware.NewPingMiddleware, middleware.NewLogMiddleware, middleware.NewServeMiddleware)
	result := chain(&request)

	c.Assert(result.Merchant, Equals, request.Merchant)
}
