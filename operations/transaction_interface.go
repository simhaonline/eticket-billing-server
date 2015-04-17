package operations

import (
	driver "database/sql/driver"
	"encoding/xml"
	"fmt"
	"time"
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

type OperationError struct {
	Message string
	Code    string
}

func (t *OperationError) String() string {
	return fmt.Sprintf("%v: %v", t.Code, t.Message)
}

func (t *OperationError) Error() string {
	return t.String()
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
