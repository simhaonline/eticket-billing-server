package operations

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "fmt"
)

type BudgetTestSuite struct {
    suite.Suite
}

func (suite *BudgetTestSuite) TestCalculate() {
    r1 := NewTransaction(fmt.Sprintf(xmlData, 101, 20200))
    r1.Save()
    r2 := NewTransaction(fmt.Sprintf(xmlData, 102, 33000))
    r2.Save()

    b := Budget{Merchant: "11", Amount: 0}
    result, _ := b.Calculate()

    suite.Equal(53200, result, "Should calculate sum of amount")
    suite.Equal(53200, b.Amount, "Should calculate sum of amount")
}

func (suite *BudgetTestSuite) TestXmlResponse() {
    budget := Budget{Merchant: "1", Amount: 123}
    answer := []byte(`<response type="budget"><merchant>1</merchant><amount>123</amount></response>`)
    answer = append(answer, '\n')
    suite.Equal(answer, budget.XmlResponse(), "Wrong xml answer")
}

func TestBudgetSuite(t *testing.T) {
    suite.Run(t, new(BudgetTestSuite))
}
