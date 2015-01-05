package server

import (
    "net"
    "encoding/xml"
)

type Request struct {
    Conn net.Conn
    XMLName xml.Name `xml:"request"`
    Merchant string `xml:"merchant"`
    OperationType string `xml:"type,attr"`
    XmlBody string
}

func NewRequest(xmlData string) *Request {
    r := &Request{}

    err := xml.Unmarshal([]byte(xmlData), r)
    if err != nil { panic(err) }

    r.XmlBody = xmlData

    return r
}
