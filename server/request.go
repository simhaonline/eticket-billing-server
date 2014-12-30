package server

import (
    "net"
)

type Request struct {
    connection net.Conn
    body string
}
