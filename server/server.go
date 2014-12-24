package server

import (
	"log"
	"net"
	"os"
	"strings"
	"strconv"
)

var (
	Info *log.Logger
	Error *log.Logger
	listOfHanldlers = make([]*Handler, 0)
)

type Handler struct {
	merchant int
	handler func(c net.Conn, b string)
}

func findOrCreateHandler(merchant int) *Handler {
	pos := -1
	for ind, elem := range listOfHanldlers {
		if elem.merchant == merchant {
			pos = ind
		}
	}
	Info.Printf("found pos: %v for merchant %v", pos, merchant)

	if pos == -1 {
		h := createHandler(merchant)
		listOfHanldlers = append(listOfHanldlers, h)
		Info.Println("create new handler")
		pos = 0
	}

	return listOfHanldlers[pos]
}

func init() {
	Info = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}

func Serve() {
	Info.Println("Ready")
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		Error.Println(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			Error.Println(err)
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			Error.Println("Error reading:", err.Error())
		}

		input := string(buf)
		Info.Println(input)
		params := strings.Split(input, "###")

		mer, err := strconv.Atoi(params[0])

		handler := findOrCreateHandler(mer)

		go handler.handler(conn, params[1])
	}
}

func createHandler(num int) *Handler {
	h := &Handler{merchant: num, handler: func(con net.Conn, body string) {
		Info.Println(body)
		con.Write([]byte("Confirm\n"))

		con.Close()
	}}

	return h
}
