package operations

import (
    "testing"
    "reflect"
    "github.com/stretchr/testify/suite"
    "time"
    "fmt"
)

type TransactionTestSuite struct {
    suite.Suite
}

func (suite TransactionTestSuite) SetupTest() {
    conn := NewConnection()
    defer conn.Close()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}

func (suite TransactionTestSuite) TearDownTest() {
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

func (suite *TransactionTestSuite) TestNewTransaction() {
    record := NewTransaction(fmt.Sprintf(xmlData, 101, 12387))

    suite.Equal("*operations.Transaction", reflect.TypeOf(record).String(), "NewTransaction should return new record composed from xml")
    suite.Equal("11", record.Merchant)
    suite.Equal("101", record.OperationIdent)
    suite.Equal("Charge", record.Description)
    suite.Equal(12387, record.Amount)

    ct := customTime{time.Date(2014, time.October, 1, 20, 13, 56, 0, time.UTC)}
    suite.Equal(ct, record.OperationCreatedAt)
}

func (suite *TransactionTestSuite) TestSave() {
    initialValue := countRows()
    record := NewTransaction(fmt.Sprintf(xmlData, 101, 100))
    record.Save()
    finishValue := countRows()
    suite.Equal(1, (finishValue - initialValue), "Count of records should be changed by 1")
}

func (suite *TransactionTestSuite) TestNotEnoughMoney() {
    r := NewTransaction(fmt.Sprintf(xmlData, 101, 100))
    r.Save()
    r = NewTransaction(fmt.Sprintf(xmlData, 102, -200))
    _, err := r.Save()

    error := err.(*TransactionError)

    suite.NotNil(err)
    suite.Equal("not_enough_money", error.Code, "Returns error's code")
    suite.NotNil(error.Message)
}

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
