package operations

import (
    "testing"
    "reflect"
    "github.com/stretchr/testify/suite"
    "time"
    "database/sql"
    "fmt"
)

type TransactionTestSuite struct {
    suite.Suite
}

func (suite TransactionTestSuite) TearDownTest() {
    conn := NewConnection()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}

var conn *sql.DB = NewConnection()

func countRows() uint64 {
    var cnt uint64
    ok := conn.QueryRow("select count(*) as cnt from operations").Scan(&cnt)
    if ok != nil { panic(ok) }
    return cnt
}

var xmlData string = `
<operation>
  <merchant>11</merchant>
  <operation_ident>%v</operation_ident>
  <description>Charge</description>
  <operation_created_at>2014-10-01 20:13:56</operation_created_at>
  <amount>%v</amount>
</operation>`

func (suite *TransactionTestSuite) TestNewRecord() {
    record := NewRecord(fmt.Sprintf(xmlData, 101, -12387))

    suite.Equal("*operations.Transaction", reflect.TypeOf(record).String(), "NewRecord should return new record composed from xml")
    suite.Equal("11", record.Merchant)
    suite.Equal("101", record.OperationIdent)
    suite.Equal("Charge", record.Description)
    suite.Equal(-12387, record.Amount)

    ct := customTime{time.Date(2014, time.October, 1, 20, 13, 56, 0, time.UTC)}
    suite.Equal(ct, record.OperationCreatedAt)
}

func (suite *TransactionTestSuite) TestSave() {
    initialValue := countRows()
    record := NewRecord(fmt.Sprintf(xmlData, 101, 100))
    record.Save()
    finishValue := countRows()
    suite.Equal(1, (finishValue - initialValue), "Count of records should be changed by 1")
}

func TestTransactionSuite(t *testing.T) {
    suite.Run(t, new(TransactionTestSuite))
}