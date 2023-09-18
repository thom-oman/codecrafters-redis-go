package main

import (
	"bufio"
	"fmt"
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

	conn, err := l.Accept()
	defer conn.Close()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	for {
		handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var resp []byte
	n, err := reader.Read(resp)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
	}

	if n > 0 {
		fmt.Printf("%v", resp)
	}

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Failed to write")
		os.Exit(1)
	}
}
