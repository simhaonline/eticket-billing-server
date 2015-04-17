package server

import (
	"encoding/xml"
)

type Budget struct {
	XMLName       xml.Name `xml:"response"`
	OperationType string   `xml:"type,attr"`
	Merchant      string   `xml:"merchant"`
	Amount        int64    `xml:"amount"`

	Db *DbConnection
}

func NewBudget(data string, merchant string, db *DbConnection) *Budget {
	return &Budget{Merchant: merchant, Db: db}
}

func (b Budget) XmlResponse() string {
	b.OperationType = "budget"
	output, _ := xml.Marshal(b)
	output = append(output, '\n')
	return string(output)
}

func (b *Budget) Calculate() (int64, error) {
	var transactions []Transaction

	b.Db.Db.Where("merchant_id = ?", b.Merchant).Find(&transactions)
	sum := calculateSum(transactions)

	return sum, nil
}

func calculateSum(transactions []Transaction) int64 {
	var sum int64
	for _, transaction := range transactions {
		sum = sum + transaction.Amount
	}
	return sum
}
