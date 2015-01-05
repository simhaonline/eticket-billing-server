package operations

import (
    "time"
    "encoding/xml"
    "fmt"
    _ "github.com/lib/pq"
//    "database/sql"
    driver "database/sql/driver"
)

const (
    TimeFormat = "2006-01-02 15:04:05"
)

type customTime struct {
    time.Time
}

// SEE http://play.golang.org/p/EFXZNsjE4a and
// http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields
func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var v string
    d.DecodeElement(&v, &start)
    parse, _ := time.Parse(TimeFormat, v)
    *c = customTime{parse}
    return nil
}

func (c customTime) Value() (driver.Value, error) {
    return c.Time, nil
}

type Transaction struct {
    XMLName xml.Name `xml:"operation"`
    Merchant string `xml:"merchant"`
    OperationIdent string `xml:"operation_ident"`
    Description string `xml:"description"`
    Amount int `xml:"amount"`

    OperationCreatedAt customTime `xml:"operation_created_at"`
    OriginXml string
}

func NewRecord(data string) *Transaction {
    r := Transaction{}
    err := xml.Unmarshal([]byte(data), &r)

    if err != nil {
        fmt.Printf("error: %v", err)
        // TODO pass error
        return nil
    }

    r.OriginXml = data

    return &r
}

func (r *Transaction) Save() (uint64, error) {
    var id uint64

    conn := NewConnection()
    ok := conn.QueryRow(`INSERT INTO operations (merchant_id, operation_ident, description, amount, operation_created_at, xml_data)
                         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
        r.Merchant, r.OperationIdent, r.Description, r.Amount, r.OperationCreatedAt, r.OriginXml).Scan(&id)

    if ok != nil { panic(ok) }

    return id, nil
}
