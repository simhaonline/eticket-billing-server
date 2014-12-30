package billing

import (
    "testing"
    "reflect"
    "github.com/stretchr/testify/assert"
    "time"
//	"fmt"
)

func TestNewRecord(t *testing.T) {
    assert := assert.New(t)

    xml := `<Operation><Merchant>11</Merchant><Description>Charge</Description><OperationCreatedAt>2014-10-01 20:13:56</OperationCreatedAt><Amount>-12387</Amount></Operation>`

    record := NewRecord(xml)

    assert.Equal("*billing.Record", reflect.TypeOf(record).String(), "NewRecord should return new record composed from xml")
    assert.Equal(11, record.Merchant)
    assert.Equal("Charge", record.Description)
    assert.Equal(-12387, record.Amount)

    ct := customTime{time.Date(2014, time.October, 1, 20, 13, 56, 0, time.UTC)}
    assert.Equal(ct, record.OperationCreatedAt)
}
