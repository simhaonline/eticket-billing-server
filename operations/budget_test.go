package operations

import (
    "eticket-billing-server/config"
    "testing"
    "fmt"
    . "gopkg.in/check.v1"
)

func TestBudget(t *testing.T) { TestingT(t) }

type BudgetSuite struct{}

var _ = Suite(&BudgetSuite{})

func (s *BudgetSuite) SetUpSuite(c *C) {
    config := config.NewConfig("test", "../config.gcfg")
    SetupConnections(config)
}

func (s *BudgetSuite) SetUpTest(c *C) {
    conn := NewConnection()
    defer conn.Close()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}

func (s *BudgetSuite) TearDownTest(c *C) {
    conn := NewConnection()
    defer conn.Close()
    _, ok := conn.Exec("truncate table operations")
    if ok != nil { panic(ok) }
}

func (s *BudgetSuite) TestCalculate(c *C) {
    r1 := NewTransaction(fmt.Sprintf(xmlData, 101, 20200))
    r1.Save()
    r2 := NewTransaction(fmt.Sprintf(xmlData, 102, 33000))
    r2.Save()

    b := Budget{Merchant: "11", Amount: 0}
    result, _ := b.Calculate()

    c.Assert(result, Equals, int64(53200))
    c.Assert(b.Amount, Equals, int64(53200))
}

func (s *BudgetSuite) TestXmlResponse(c *C) {
    budget := Budget{Merchant: "1", Amount: 123}
    answer := []byte(`<response type="budget"><merchant>1</merchant><amount>123</amount></response>`)
    answer = append(answer, '\n')
    c.Assert(budget.XmlResponse(), Equals, string(answer))
}
