package server

import (
    "testing"
    "encoding/xml"
    "github.com/stretchr/testify/assert"
)

var xmlData string = `
<request type="budget">
  <merchant>11</merchant>
  <operation_ident>%v</operation_ident>
  <description>Charge</description>
  <operation_created_at>2014-10-01 20:13:56</operation_created_at>
  <amount>%v</amount>
</request>`

func TestRequest(t *testing.T) {
    assert := assert.New(t)

    r := Request{}
    _ = xml.Unmarshal([]byte(xmlData), &r)
    assert.Equal("11", r.Merchant)
    assert.Equal("budget", r.OperationType)
}
