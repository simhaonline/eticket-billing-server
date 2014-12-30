package billing

import (
    "time"
    "encoding/xml"
    "fmt"
    _ "github.com/lib/pq"
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

type Record struct {
    XMLName xml.Name `xml:"Operation"`
    Merchant string `xml:"Merchant"`
    OperationIdent string `xml:OperationIdent`
    Description string `xml:"Description"`
    Amount int `xml:"Amount"`

    OperationCreatedAt customTime `xml:"OperationCreatedAt"`
    OriginXml string
}

type Budget struct {
    Merchant int
    Amount int
}

func NewRecord(data string) *Record {
    r := Record{}
    err := xml.Unmarshal([]byte(data), &r)

    if err != nil {
        fmt.Printf("error: %v", err)
        // TODO pass error
        return nil
    }

    r.OriginXml = data

    return &r
}

func (r *Record) Save() (uint64, error) {
    var id uint64

    conn := NewConnection()
    ok := conn.QueryRow(`INSERT INTO operations (merchant_id, operation_ident, description, amount, operation_created_at, xml_data)
                         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
        r.Merchant, r.OperationIdent, r.Description, r.Amount, r.OperationCreatedAt, r.OriginXml).Scan(&id)

    if ok != nil { panic(ok) }

    return id, nil
}

func (b *Budget) Calculate() (int, error) {
    conn := NewConnection()
    ok := conn.QueryRow(`select sum(amount) from operations where merchant_id = $1`, b.Merchant).Scan(&b.Amount)

    if ok != nil { panic(ok) }

    return b.Amount, nil
}
