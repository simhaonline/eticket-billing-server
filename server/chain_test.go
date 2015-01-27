package server

import (
    "testing"
    "eticket-billing-server/request"
    . "gopkg.in/check.v1"
)

func TestChain(t *testing.T) { TestingT(t) }

type ChainSuite struct{
    CheckArray []string
}

var _ = Suite(&ChainSuite{CheckArray: []string{}})

var request request.Request

func (s *ChainSuite) SetUpSuite(c *C) {

}

func (s *ChainSuite) TestNewChain(c *C) {
    request = request.Request{Merchant: "m10"}
    chain := NewChain(NewPingMiddleware, NewLogMiddleware, NewServeMiddleware)
    result := chain(&request)

    c.Assert(result.Merchant, Equals, request.Merchant)
}
