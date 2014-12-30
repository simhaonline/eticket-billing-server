package billing

import (
    "testing"
    "reflect"
    "github.com/stretchr/testify/suite"
    "time"
    "database/sql"
    "fmt"
)

type TestSuite struct {
    suite.Suite
}

func (suite TestSuite) TearDownTest() {
    conn := NewConnection()
    _, ok := conn.Exec("delete from operations")
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
<Operation>
  <Merchant>11</Merchant>
  <OperationIdent>%v</OperationIdent>
  <Description>Charge</Description>
  <OperationCreatedAt>2014-10-01 20:13:56</OperationCreatedAt>
  <Amount>%v</Amount>
</Operation>`

func (suite *TestSuite) TestNewRecord() {
    record := NewRecord(fmt.Sprintf(xmlData, 101, -12387))

    suite.Equal("*billing.Record", reflect.TypeOf(record).String(), "NewRecord should return new record composed from xml")
    suite.Equal("11", record.Merchant)
    suite.Equal("101", record.OperationIdent)
    suite.Equal("Charge", record.Description)
    suite.Equal(-12387, record.Amount)

    ct := customTime{time.Date(2014, time.October, 1, 20, 13, 56, 0, time.UTC)}
    suite.Equal(ct, record.OperationCreatedAt)
}

func (suite *TestSuite) TestSave() {
    initialValue := countRows()
    record := NewRecord(fmt.Sprintf(xmlData, 101, 100))
    record.Save()
    finishValue := countRows()
    suite.Equal(1, (finishValue - initialValue), "Count of records should be changed by 1")
}

func (suite *TestSuite) TestCalculate() {
    r1 := NewRecord(fmt.Sprintf(xmlData, 101, 20200))
    r1.Save()
    r2 := NewRecord(fmt.Sprintf(xmlData, 102, 33000))
    r2.Save()

    b := Budget{11, 0}
    result, _ := b.Calculate()

    suite.Equal(53200, result, "Should calculate sum of amount")
    suite.Equal(53200, b.Amount, "Should calculate sum of amount")
}

func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(TestSuite))
}
