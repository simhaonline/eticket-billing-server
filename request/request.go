package request

import (
    "net"
    "encoding/xml"
    "runtime"
    "fmt"
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

func (req *Request) Perform(fnc func(r *Request) string) *Request {
    defer func() {
        if err := recover(); err != nil {
            trace := make([]byte, 1024)
            count := runtime.Stack(trace, true)
            fmt.Printf("Recover from panic: %s\n", err)
            fmt.Printf("Stack of %d bytes: %s\n", count, trace)
        }
    }()
    defer req.Conn.Close()
    xmlString := fnc(req)
    req.Conn.Write([]byte(xmlString))
    return req
}
