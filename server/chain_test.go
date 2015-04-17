package server

import (
	"eticket-billing-server/config"
	. "gopkg.in/check.v1"
	"testing"
	"net"
	"time"
)

func TestChain(t *testing.T) { TestingT(t) }

type ChainSuite struct {
	CheckArray []string
	db *DbConnection
	mapping  *PerformerFnMapping
}

var _ = Suite(&ChainSuite{CheckArray: []string{}})

type BarePerformer struct {
	Request *Request
	Db      *DbConnection
}

func (p *BarePerformer) Serve() *Request {
	return p.Request.Perform(func(req *Request) string {
		return ""
	})
}

func NewBarePerformer(request *Request, connection *DbConnection) performerType {
	b := BarePerformer{Request: request, Db: connection}
	return performerType(&b)
}

type FakeConn struct {}

func (f FakeConn) Close() error { return nil }
func (f FakeConn) LocalAddr() net.Addr { return nil }
func (f FakeConn) RemoteAddr() net.Addr { return nil }
func (f FakeConn) Read(b []byte) (n int, err error) { return 0, nil }
func (f FakeConn) Write(b []byte) (n int, err error) { return 0, nil }
func (f FakeConn) SetDeadline(t time.Time) error { return nil}
func (f FakeConn) SetReadDeadline(t time.Time) error { return nil }
func (f FakeConn) SetWriteDeadline(t time.Time) error { return nil }

func (s *ChainSuite) SetUpSuite(c *C) {
	m := make(PerformerFnMapping)
	m["bare"] = NewBarePerformer
	s.mapping = &m


	config := config.NewConfig("test", "../config.gcfg")
	s.db = NewConnection(config)
}

func (s *ChainSuite) TestNewChain(c *C) {
	request := Request{Merchant: "m10", OperationType: "bare", Conn: FakeConn{}}
	context := Context{Request: &request, Db: s.db , PerformersMapping: s.mapping}
	chain := NewChain(NewPingMiddleware, NewLogMiddleware, NewServeMiddleware)
	result := chain(&context)

	c.Assert(result.Request.Merchant, Equals, request.Merchant)
}
