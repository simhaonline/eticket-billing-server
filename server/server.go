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
	listOfWorkers = make([]*Worker, 0)
)

type Request struct {
	connection net.Conn
	body string
}

type Worker struct {
	merchant int
	inputChan chan *Request
	quitChan chan bool
}

func (w Worker) Serve() {
	Info.Printf("New worker for merchant %v is spawned", w.merchant)
	var req *Request
	for {
		select {
		case req = <- w.inputChan:
			Info.Printf("Worker of merchant %v received string '%v'", w.merchant, req.body)
			req.connection.Write([]byte("Confirm\n"))
			req.connection.Close()
		case <- w.quitChan:
			Info.Println("Wroker for %v quitting", w.merchant)
			return
		}
	}
}

func init() {
	Info = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}

func findOrCreateWorker(merchant int) *Worker {
	pos := getPosition(merchant)

	Info.Printf("found pos: %v for merchant %v", pos, merchant)

	if -1 == pos {
		h := createWorker(merchant)
		listOfWorkers = append(listOfWorkers, h)

		pos = getPosition(merchant)

		Info.Printf("Spawning new worker for merchant %v", merchant)
		go listOfWorkers[pos].Serve()
	}

	return listOfWorkers[pos]
}

func getPosition(merchant int) int {
	pos := -1
	for ind, elem := range listOfWorkers {
		if elem.merchant == merchant {
			pos = ind
		}
	}
	return pos
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

		worker := findOrCreateWorker(mer)
		worker.inputChan <- &Request{conn, params[1]}
	}
}

func createWorker(merchant int) *Worker {
	h := &Worker{ merchant, make(chan *Request), make(chan bool) }
	return h
}
