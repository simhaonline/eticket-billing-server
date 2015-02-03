// +build helpprog

package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	counter := 0

	strCheck := "<request type='budget'><merchant>1</merchant></request>"
	strTransaction := "<request type='transaction'><merchant>1</merchant><operation_ident>aw%v</operation_ident><description>hello</description><amount>1</amount><operation_created_at>2014-10-01 20:13:56</operation_created_at></request>"

	for r := 0; r < 100; r++ {
		go func() {
			servAddr := "localhost:2000"
			tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
			if err != nil {
				println("ResolveTCPAddr failed:", err.Error())
				os.Exit(1)
			}

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			defer conn.Close()

			if err != nil {
				println("Dial failed:", err.Error())
				os.Exit(1)
			}

			uuid := uuid.NewUUID()

			req := fmt.Sprintf(strTransaction, uuid.String()[:10])
			counter++

			_, err = conn.Write([]byte(req))
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}

			fmt.Println("write to server = ", req)

			runtime.Gosched()
		}()
	}

	time.Sleep(time.Second)

	servAddr := "localhost:2000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	req := fmt.Sprintf(strCheck, counter)
	counter++

	_, err = conn.Write([]byte(req))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("write to server = ", req)

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("reply from server=", string(reply))

	os.Exit(0)
}
