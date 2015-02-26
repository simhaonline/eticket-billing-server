package server

import (
	"encoding/xml"
	. "gopkg.in/check.v1"
	"testing"
)

func TestRequest(t *testing.T) { TestingT(t) }

type RequestSuite struct{}

var _ = Suite(&RequestSuite{})

var xmlData string = `
<request type="budget">
  <merchant>11</merchant>
  <operation_ident>%v</operation_ident>
  <description>Charge</description>
  <operation_created_at>2014-10-01 20:13:56</operation_created_at>
  <amount>%v</amount>
</request>`

func (s *RequestSuite) TestRequest(c *C) {
	r := Request{}
	_ = xml.Unmarshal([]byte(xmlData), &r)
	c.Assert(r.Merchant, Equals, "11")
	c.Assert(r.OperationType, Equals, "budget")
}
