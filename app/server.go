package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

var (
	SEPARATOR = []byte{13, 10}
)

func main() {
	fmt.Printf("Listening on port %v\n", CONN_PORT)
	l, err := net.Listen(CONN_TYPE, "0.0.0.0:"+CONN_PORT)
	if err != nil {
		fmt.Println("Failed to bind to port " + CONN_PORT)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error parsing incoming message: ", err.Error())
			os.Exit(1)
		}
		// fmt.Println(buf)
		cmds := splitRequest(buf)
		fmt.Println(cmds)
		for i, _ := range cmds {
			fmt.Println(string(cmds[i]))
		}
		writeResponse(conn, []byte("+PONG\r\n"))
	}
}

func writeResponse(conn net.Conn, resp []byte) {
	_, err := conn.Write(resp)
	if err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}

func splitRequest(req []byte) [][]byte {
	var cmds [][]byte
	var cur []byte
	bulk := isBulkRequest(req[0])

	for i := range req {
		if i == 0 {
			continue
		}

		if req[i] == 0 {
			break
		}
		if !bulk && (req[i] == 10 || req[i] == 13) {
			fmt.Println("Cannot provide \\n or \\r")
			os.Exit(1)
		}
		cur = append(cur, req[i])

		if len(cur) >= 2 && bytes.Equal(cur[len(cur)-2:], SEPARATOR) {
			if !bulk {
				fmt.Printf("Complex command provided: %v\n", req[0])
				os.Exit(1)
			}

			cmds = append(cmds, cur[:len(cur)-2])
			cur = nil
		}
	}

	cmds = append(cmds, cur)
	return cmds
}

func isBulkRequest(b byte) bool {
	switch b {
	case 43, 45, 58, 95, 35, 44, 40:
		// simple
		return false
	case 36, 42, 33, 61, 37, 126, 62:
		return true
	default:
		return false
	}
}
