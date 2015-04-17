package operations

import (
	"eticket-billing-server/config"
	"fmt"
	gorm "github.com/jinzhu/gorm"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
	"time"
)

func TestTransaction(t *testing.T) { TestingT(t) }

type TransactionSuite struct {
	db *gorm.DB
}

var _ = Suite(&TransactionSuite{})

func (s *TransactionSuite) SetUpSuite(c *C) {
	config := config.NewConfig("test", "../config.gcfg")
	SetupConnections(config)
	s.db = NewConnection()
}

func (s *TransactionSuite) SetUpTest(c *C) {
	s.db.Exec("truncate table operations")
}

func (s *TransactionSuite) TearDownTest(c *C) {
	s.db.Exec("truncate table operations")
}

func countRows(s *TransactionSuite) uint64 {
	var count uint64
	s.db.Model(Transaction{}).Count(&count)

	return count
}

var xmlData string = `
<request type="transaction">
  <application_name>app</application_name>
  <merchant>11</merchant>
  <operation_name>charge</operation_name>
  <operation_ident>%v</operation_ident>
  <description>Charge</description>
  <operation_created_at>2014-10-01 20:13:56</operation_created_at>
  <amount>%v</amount>
</request>`

func (s *TransactionSuite) TestNewTransaction(c *C) {
	record := NewTransaction(fmt.Sprintf(xmlData, 101, 12387))

	c.Assert(reflect.TypeOf(record).String(), Equals, "*operations.Transaction")
	c.Assert(record.Merchant, Equals, "11")
	c.Assert(record.OperationIdent, Equals, "101")
	c.Assert(record.Description, Equals, "Charge")
	c.Assert(record.Amount, Equals, int64(12387))

	ct := customTime{time.Date(2014, time.October, 1, 20, 13, 56, 0, time.UTC)}
	c.Assert(record.OperationCreatedAt, Equals, ct)
}

func (s *TransactionSuite) TestSave(c *C) {
	initialValue := countRows(s)
	record := NewTransaction(fmt.Sprintf(xmlData, 101, 100))

	s.db.Exec("select count(*)")
	s.db.Create(record)
	finishValue := countRows(s)
	c.Assert(int(finishValue-initialValue), Equals, 1)
}

func (s *TransactionSuite) TestNotEnoughMoney(c *C) {
	record1 := NewTransaction(fmt.Sprintf(xmlData, 101, 100))
	s.db.Create(record1)

	record2 := NewTransaction(fmt.Sprintf(xmlData, 102, -200))
	result := s.db.Create(record2)

	error := result.Error

	c.Assert(error, NotNil)
	c.Assert(error.Error(), Equals, "Not enough money for operation")
	c.Assert(record2.Errors[0].Code, Equals, "not_enough_money")
	c.Assert(record2.Errors[0].Message, Equals, "Not enough money for operation")
}

func (s *TransactionSuite) TestDuplicationOfRecords(c *C) {
	initialValue := countRows(s)
	record1 := NewTransaction(fmt.Sprintf(xmlData, 101, 12387))
	s.db.Create(record1)

	record2 := NewTransaction(fmt.Sprintf(xmlData, 101, 12387))
	result := s.db.Create(record2)

	finishValue := countRows(s)

	c.Assert(result.Error, NotNil)

	error := result.Error

	c.Assert(error.Error(), NotNil)

	c.Assert(int(initialValue), Equals, 0)
	c.Assert(int(finishValue), Equals, 1)

	xml := record2.ErrorXmlResponse(result.Error)
	answer := []byte(`<answer type="transaction"><error><message>Duplicate key value violates unique constraint</message><code>unique_violation</code></error><application_name>app</application_name><merchant>11</merchant><operation_name>charge</operation_name><operation_ident>101</operation_ident><description>Charge</description><amount>12387</amount><operation_created_at>2014-10-01T20:13:56Z</operation_created_at></answer>`)
	answer = append(answer, '\n')
	c.Assert(xml, Equals, string(answer))
}

func (s *TransactionSuite) TestXmlResponse(c *C) {
	transaction := Transaction{Merchant: "10", OperationIdent: "asdf", Description: "Hello", Amount: 101, ApplicationName: "app", OperationName: "charge"}
	answer := []byte(`<answer type="transaction"><application_name>app</application_name><merchant>10</merchant><operation_name>charge</operation_name><operation_ident>asdf</operation_ident><description>Hello</description><amount>101</amount><operation_created_at>0001-01-01T00:00:00Z</operation_created_at></answer>`)
	answer = append(answer, '\n')
	c.Assert(transaction.XmlResponse(), Equals, string(answer))
}

func (s *TransactionSuite) TestErrorXmlResponse(c *C) {
	record := NewTransaction(fmt.Sprintf(xmlData, 101, -100))
	result := s.db.Create(record)

	xml := record.ErrorXmlResponse(result.Error)
	answer := []byte(`<answer type="transaction"><error><message>Not enough money for operation</message><code>not_enough_money</code></error><application_name>app</application_name><merchant>11</merchant><operation_name>charge</operation_name><operation_ident>101</operation_ident><description>Charge</description><amount>-100</amount><operation_created_at>2014-10-01T20:13:56Z</operation_created_at></answer>`)
	answer = append(answer, '\n')
	c.Assert(xml, Equals, string(answer))
}
