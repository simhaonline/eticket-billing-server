package operations

import (
	"encoding/xml"
)

type Budget struct {
	XMLName       xml.Name `xml:"response"`
	OperationType string   `xml:"type,attr"`
	Merchant      string   `xml:"merchant"`
	Amount        int64    `xml:"amount"`
}

func (b Budget) XmlResponse() string {
	b.OperationType = "budget"
	output, _ := xml.Marshal(b)
	output = append(output, '\n')
	return string(output)
}

func (b *Budget) Calculate() (int64, error) {
	conn := NewConnection()
	defer conn.Close()

	var transactions []Transaction
	db := GetDB()
	db.Where("merchant_id = ?", b.Merchant).Find(&transactions)
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
