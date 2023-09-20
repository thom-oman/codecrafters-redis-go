package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
		req, _ := buildRequest(buf)
		switch strings.ToLower(req.Command()) {
		case "ping":
			writeResponse(conn, []byte("+PONG\r\n"))
		case "echo":
			writeResponse(conn, []byte(req.Args()[0]+"\r\n"))
		}
		fmt.Println(req)
		// for i, _ := range req {
		// 	fmt.Println(string(req[i]))
		// }
	}
}

type request interface {
	Command() string
	Args() []string
	Simple() bool
	// Add([]byte)
}

type simpleStringsRequest struct {
	cmd string
}

func (r simpleStringsRequest) Command() string {
	return r.cmd
}

func (r simpleStringsRequest) Args() []string {
	return []string{}
}

func (r simpleStringsRequest) Simple() bool {
	return true
}

type arrayRequest struct {
	cmd  string
	args []string
	size int
}

func (r arrayRequest) Command() string {
	return r.cmd
}

func (r arrayRequest) Args() []string {
	return r.args
}

func (r arrayRequest) Simple() bool {
	return false
}

func writeResponse(conn net.Conn, resp []byte) {
	_, err := conn.Write(resp)
	if err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}

func buildRequest(req []byte) (request, error) {
	fmt.Println(req)
	d := split(req[1:])
	fmt.Println(d)

	var r request

	switch req[0] {
	case 43: // +
		// Simple strings	Simple	+
		if len(d) != 1 {
			fmt.Printf("Not sure we can have multiple arguments to a simple strings request")
			os.Exit(1)
		}
		r = &simpleStringsRequest{cmd: string(d[0])}
	case 42: // *
		// Arrays	Aggregate	*
		l, err := strconv.Atoi(string(d[0]))
		if err != nil {
			fmt.Printf("Could not convert %v to int: %v", d[0], err.Error())
			os.Exit(1)
		}

		if l != len(d)-1 {
			fmt.Printf("Expected %v elements but %v provided", l, len(d)-1)
			os.Exit(1)
		}
		var args []string

		for i := 2; i < len(d); i++ {
			args = append(args, string(d[i]))
		}

		r = &arrayRequest{
			cmd:  string(d[1]),
			args: args,
			size: l,
		}
	}
	if r == nil {
		return nil, errors.New(fmt.Sprintf("Invalid first token %v", req[0]))
	}
	return r, nil
}

func split(req []byte) [][]byte {
	var cmds [][]byte
	var cur []byte

	for i := range req {
		if req[i] == 0 {
			break
		}
		// fmt.Printf("Appending %v to %v. Complete: %v\n", req[i], cur, cmds)
		cur = append(cur, req[i])

		if len(cur) >= 2 && bytes.Equal(cur[len(cur)-2:], SEPARATOR) {
			if len(cur[:len(cur)-2]) > 0 {
				cmds = append(cmds, cur[:len(cur)-2])
			}
			cur = nil
		}
	}
	if len(cur) != 0 {
		// Should this ever be the case?
		cmds = append(cmds, cur)
	}
	return cmds
}
