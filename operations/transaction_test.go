package operations

import (
    "eticket-billing-server/config"
    "testing"
    "reflect"
    "time"
    "fmt"
    . "gopkg.in/check.v1"
)

func TestTransaction(t *testing.T) { TestingT(t) }

type TransactionSuite struct{}

var _ = Suite(&TransactionSuite{})

func (s *TransactionSuite) SetUpSuite(c *C) {
    config := config.NewConfig("test", "../config.gcfg")
    SetupConnections(config)
}

func (s *TransactionSuite) SetUpTest(c *C) {
    conn := NewConnection()
    defer conn.Close()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}

func (s *TransactionSuite) TearDownTest(c *C) {
    conn := NewConnection()
    defer conn.Close()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}


func countRows() uint64 {
    var cnt uint64
    conn := NewConnection()
    defer conn.Close()
    ok := conn.QueryRow("select count(*) as cnt from operations").Scan(&cnt)
    if ok != nil { panic(ok) }
    return cnt
}

var xmlData string = `
<request type="transaction">
  <merchant>11</merchant>
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
    initialValue := countRows()
    record := NewTransaction(fmt.Sprintf(xmlData, 101, 100))
    record.Save()
    finishValue := countRows()
    c.Assert(int(finishValue - initialValue), Equals, 1)
}

func (s *TransactionSuite) TestNotEnoughMoney(c *C) {
    r := NewTransaction(fmt.Sprintf(xmlData, 101, 100))
    r.Save()
    r = NewTransaction(fmt.Sprintf(xmlData, 102, -200))
    _, err := r.Save()

    error := err.(*TransactionError)

    c.Assert(err, NotNil)
    c.Assert(error.Code, Equals, "not_enough_money")
    c.Assert(error.Message, NotNil)
}

/*
func (suite *TransactionTestSuite) TestDuplicationOfRecords() {
    initialValue := countRows()
    r := NewTransaction(fmt.Sprintf(xmlData, 101, 12387))
    r.Save()
    r = NewTransaction(fmt.Sprintf(xmlData, 101, 12387))
    _, err := r.Save()
    finishValue := countRows()

    error := err.(*TransactionError)

    suite.NotNil(err)
    suite.Equal("duplication", error.Code, "Returns error's code")
    suite.NotNil(error.Message)
    suite.Equal(0, int(initialValue), "First must be 1")
    suite.Equal(1, int(finishValue), "Result must be 1")
}

func (suite *TransactionTestSuite) TestXmlResponse() {
    transaction := Transaction{Merchant: "10", OperationIdent: "asdf", Description: "Hello", Amount: 101}
    answer := []byte(`<answer type="transaction"><merchant>10</merchant><operation_ident>asdf</operation_ident><description>Hello</description><amount>101</amount><operation_created_at>0001-01-01T00:00:00Z</operation_created_at></answer>`)
    answer = append(answer, '\n')
    suite.Equal(string(answer), transaction.XmlResponse(), "Wrong xml answer")
}

func (suite *TransactionTestSuite) TestErrorXmlResponse() {
    r := NewTransaction(fmt.Sprintf(xmlData, 101, -100))
    _, err := r.Save()

    xml := r.ErrorXmlResponse(err)
    answer := []byte(`<answer type="transaction"><error><message>Not enough money for operation</message><code>not_enough_money</code></error><merchant>11</merchant><operation_ident>101</operation_ident><description>Charge</description><amount>-100</amount><operation_created_at>2014-10-01T20:13:56Z</operation_created_at></answer>`)
    answer = append(answer, '\n')
    suite.Equal(string(answer), xml, "Compose xml with error")
}


func TestTransactionSuite(t *testing.T) {
    suite.Run(t, new(TransactionTestSuite))
}
*/
