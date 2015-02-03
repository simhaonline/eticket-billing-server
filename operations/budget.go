package operations

import (
	"database/sql"
	"encoding/xml"
	_ "github.com/lib/pq"
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

	var amount sql.NullInt64

	ok := conn.QueryRow("select sum(amount) from operations where merchant_id = $1", b.Merchant).Scan(&amount)
	if ok != nil {
		panic(ok)
	}

	if amount.Valid {
		b.Amount = amount.Int64
	} else {
		b.Amount = 0
	}

	return b.Amount, nil
}
