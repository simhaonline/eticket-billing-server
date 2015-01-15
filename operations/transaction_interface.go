package operations


import (
    "time"
    driver "database/sql/driver"
    "fmt"
    "encoding/xml"
)

const (
    TimeFormat = "2006-01-02 15:04:05"
)

type ITransaction interface {
    Save() (uint64, error)
    XmlResponse() string
    ErrorXmlResponse(err error) string
}

type customTime struct {
    time.Time
}

type TransactionError struct {
    Message string
    Code string
}

func (t *TransactionError) Error() string {
    return fmt.Sprintf("%v: %v", t.Code, t.Message)
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
