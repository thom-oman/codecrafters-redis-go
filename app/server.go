package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/thom-oman/codecrafters-redis-go/app/parser"
	"github.com/thom-oman/codecrafters-redis-go/app/store"
)

const (
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

// var (
// 	SEPARATOR = []byte{13, 10}
// )

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
	// defer closeConnection(conn)

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error parsing incoming message: ", err.Error())
			os.Exit(1)
		}
		req := parser.NewRequest()
		fmt.Printf("Parsing %v\n", string(buf))
		req.AddTokens(buf)
		req.Parse()
		args := req.Args()
		cmd, params := args[0], args[1:]

		switch strings.ToLower(cmd) {
		case "ping":
			writeResponse(conn, []byte("+PONG\r\n"))
		case "echo":
			writeResponse(conn, []byte("+"+params[0]+"\r\n"))
		case "set":
			if len(params) <= 2 {
				fmt.Println("Must supply 2 arguments to SET")
			}
			key, value, xp := params[0], params[1], params[2:]
			px := -1
			if len(xp) != 0 {

			}
			fmt.Printf("SET Key: %v, Value: %v\n", key, value)
			_ = store.Set(key, value, px)
			writeResponse(conn, []byte("+OK\r\n"))
		case "get":
			if len(params) != 1 {
				fmt.Println("Must supply 1 arguments to GET")
			}
			key := params[0]
			value, _ := store.Get(key)
			resp := fmt.Sprintf("$%v\r\n%v\r\n", len(value), value)
			writeResponse(conn, []byte(resp))
		}
	}
}

func closeConnection(conn net.Conn) {
	fmt.Println("Closing connection")
	conn.Close()
}

func writeResponse(conn net.Conn, resp []byte) {
	_, err := conn.Write(resp)
	if err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}

// 	var r request

// 	switch req[0] {
// 	case 43: // +
// 		// Simple strings	Simple	+
// 		if len(d) != 1 {
// 			fmt.Printf("Not sure we can have multiple arguments to a simple strings request")
// 			os.Exit(1)
// 		}
// 		r = &simpleStringsRequest{cmd: string(d[0])}
// 	case 42: // *
// 		// Arrays	Aggregate	*
// 		l, err := strconv.Atoi(string(d[0]))
// 		if err != nil {
// 			// TODO return err
// 			fmt.Printf("Could not convert %v to int: %v", d[0], err.Error())
// 			os.Exit(1)
// 		}

// 		if l != len(d)-1 {
// 			// TODO return err
// 			fmt.Printf("Expected %v elements but %v provided", l, len(d)-1)
// 			os.Exit(1)
// 		}
// 		var args []string

// 		for i := 2; i < len(d); i++ {
// 			args = append(args, string(d[i]))
// 		}

// 		r = &arrayRequest{
// 			cmd:  string(d[1]),
// 			args: args,
// 			size: l,
// 		}
// 	}

// 	if r == nil {
// 		return nil, errors.New(fmt.Sprintf("Invalid first token %v", req[0]))
// 	}
// 	return r, nil
// }

// func split(req []byte) [][]byte {
// 	var cmds [][]byte
// 	var cur []byte

// 	for i := range req {
// 		if req[i] == 0 {
// 			break
// 		}
// 		// fmt.Printf("Appending %v to %v. Complete: %v\n", req[i], cur, cmds)
// 		cur = append(cur, req[i])

// 		if len(cur) >= 2 && bytes.Equal(cur[len(cur)-2:], SEPARATOR) {
// 			if len(cur[:len(cur)-2]) > 0 {
// 				cmds = append(cmds, cur[:len(cur)-2])
// 			}
// 			cur = nil
// 		}
// 	}
// 	if len(cur) != 0 {
// 		// Should this ever be the case?
// 		cmds = append(cmds, cur)
// 	}
// 	return cmds
// }
