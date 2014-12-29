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
	IncomeRequestsLog *log.Logger
	listOfWorkers = make([]*Worker, 0)
)


func init() {
	Info = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)

	f, ok := os.OpenFile("/tmp/income_requests.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if ok != nil { panic(ok) }

//	defer f.Close() TODO how to close it? Where?

	IncomeRequestsLog = log.New(f, "", log.Ldate|log.Ltime)
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

		// TODO should read all string
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)

		if err != nil {
			Error.Println("Error reading:", err.Error())
		}

		input := string(buf)
		Info.Println(input)
		params := strings.Split(input, "###")
		merchant, err := strconv.Atoi(params[0])

		IncomeRequestsLog.Println(input)

		worker := GetWorkerForMerchant(merchant)
		worker.inputChan <- &Request{conn, params[1]}
	}
}
