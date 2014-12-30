package billing

import (
    "time"
    "encoding/xml"
    "fmt"
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

type Record struct {
    XMLName xml.Name `xml:"Operation"`
    Merchant int `xml:"Merchant"`
    Description string `xml:"Description"`
    Amount int `xml:"Amount"`

    OperationCreatedAt customTime `xml:"OperationCreatedAt"`
    OriginXml string
}

func NewRecord(data string) *Record {
    r := Record{}
    err := xml.Unmarshal([]byte(data), &r)

    if err != nil {
        fmt.Printf("error: %v", err)
        // TODO pass error
        return nil
    }

    return &r
}

func (r *Record) Save() error {

}
