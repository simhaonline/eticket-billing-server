package server

import (
	"eticket-billing-server/config"
	"fmt"
	gorm "github.com/jinzhu/gorm"
	. "gopkg.in/check.v1"
	"testing"
)

func TestBudget(t *testing.T) { TestingT(t) }

type BudgetSuite struct {
	db *gorm.DB
}

var _ = Suite(&BudgetSuite{})

func (s *BudgetSuite) SetUpSuite(c *C) {
	config := config.NewConfig("test", "../config.gcfg")
	SetupConnections(config)
	s.db = NewConnection()
}

func (s *BudgetSuite) SetUpTest(c *C) {
	s.db.Exec("truncate table operations")
}

func (s *BudgetSuite) TearDownTest(c *C) {
	s.db.Exec("truncate table operations")
}

func (s *BudgetSuite) TestCalculate(c *C) {
	r1 := NewTransaction(fmt.Sprintf(xmlData, 101, 20200), s.db)
	s.db.Create(r1)
	r2 := NewTransaction(fmt.Sprintf(xmlData, 102, 33000))
	s.db.Create(r2)

}

func (s *BudgetSuite) TestXmlResponse(c *C) {
	budget := Budget{Merchant: "1", Amount: 123}
	answer := []byte(`<response type="budget"><merchant>1</merchant><amount>123</amount></response>`)
	answer = append(answer, '\n')
	c.Assert(budget.XmlResponse(), Equals, string(answer))
}
