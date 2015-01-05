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
    r1 := NewRecord(fmt.Sprintf(xmlData, 101, 20200))
    r1.Save()
    r2 := NewRecord(fmt.Sprintf(xmlData, 102, 33000))
    r2.Save()

    b := Budget{"11", 0}
    result, _ := b.Calculate()

    suite.Equal(53200, result, "Should calculate sum of amount")
    suite.Equal(53200, b.Amount, "Should calculate sum of amount")
}

func TestBudgetSuite(t *testing.T) {
    suite.Run(t, new(BudgetTestSuite))
}
