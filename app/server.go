package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func main() {
	fmt.Printf("Listening on port %v\n", CONN_PORT)
	l, err := net.Listen(CONN_TYPE, "0.0.0.0:"+CONN_PORT)
	if err != nil {
		fmt.Println("Failed to bind to port " + CONN_PORT)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		fmt.Println("Received request")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	l, err := reader.ReadBytes('\n')
	fmt.Printf("Received %v", l)

	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading:", err.Error())
		}
		conn.Close()
	}
	conn.Write([]byte("+PONG\r\n"))
	conn.Close()
}
