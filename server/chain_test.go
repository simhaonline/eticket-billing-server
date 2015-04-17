package server

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestChain(t *testing.T) { TestingT(t) }

type ChainSuite struct {
	CheckArray []string
}

var _ = Suite(&ChainSuite{CheckArray: []string{}})

func BarePerformere(req *Request) *Request {
	return req
}

func (s *ChainSuite) SetUpSuite(c *C) {
	mapping := make(PerformerFnMapping)
	mapping["bare"] = BarePerformere
	SetupMapping(mapping)
}

func (s *ChainSuite) TestNewChain(c *C) {
	request := Request{Merchant: "m10", OperationType: "bare"}
	chain := NewChain(NewPingMiddleware, NewLogMiddleware, NewServeMiddleware)
	result := chain(&request)

	c.Assert(result.Merchant, Equals, request.Merchant)
}
