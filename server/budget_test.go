package server

import (
	"eticket-billing-server/config"
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

func TestBudget(t *testing.T) { TestingT(t) }

type BudgetSuite struct {
	db *DbConnection
}

var _ = Suite(&BudgetSuite{})

func (s *BudgetSuite) SetUpSuite(c *C) {
	config := config.NewConfig("test", "../config.gcfg")
	s.db = NewConnection(config)
}

func (s *BudgetSuite) SetUpTest(c *C) {
	s.db.Db.Exec("truncate table operations")
}

func (s *BudgetSuite) TearDownTest(c *C) {
	s.db.Db.Exec("truncate table operations")
}

func (s *BudgetSuite) TestCalculate(c *C) {
	r1 := NewTransaction(fmt.Sprintf(transactionXmlData, 101, 20200), s.db)
	s.db.Db.Create(r1)
	r2 := NewTransaction(fmt.Sprintf(transactionXmlData, 102, 33000), s.db)
	s.db.Db.Create(r2)

	b := Budget{Merchant: "11", Amount: 0, Db: s.db}
	result, _ := b.Calculate()

	c.Assert(result, Equals, int64(53200))
	c.Assert(b.Amount, Equals, int64(53200))
}

func (s *BudgetSuite) TestXmlResponse(c *C) {
	// TODO check after save to database. It could return external columns
	budget := Budget{Merchant: "1", Amount: 123}
	answer := []byte(`<response type="budget"><merchant>1</merchant><amount>123</amount></response>`)
	answer = append(answer, '\n')
	c.Assert(budget.XmlResponse(), Equals, string(answer))
}
