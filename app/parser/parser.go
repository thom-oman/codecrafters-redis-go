package parser

// Parses RESP: https://redis.io/docs/reference/protocol-spec/

import (
	"fmt"
	"os"
	"strconv"
)

type request struct {
	tokens []byte
	parsed bool
	cur    int
	simple bool
	args []string
}

func NewRequest() *request {
	return &request{}
}

func (r request) TokenCount() int {
	return len(r.tokens)
}

func (r request) Parsed() bool {
	return r.parsed
}

func (r *request) AddTokens(ts []byte) {
	for i := range ts {
		r.AddToken(ts[i])
	}
}

func (r *request) AddToken(t byte) {
	if t != 0 {
		r.tokens = append(r.tokens, t)
	}
}

func (r *request) Args() []string {
	return r.args
}

func (r *request) Parse() error {
	if r.cur != 0 || r.Parsed() {
		fmt.Println("Requests already parsed")
		os.Exit(1)
	}

	parsedHeader := false
	var argLength int
	var cur []byte

	for i := 0; i < len(r.tokens); i++ {
		if r.tokens[i] == 13 && r.tokens[i+1] == 10 {
			if len(cur) > 0 {
				r.args = append(r.args, string(cur))
				cur = nil
			}
			i++
			continue
		}

		switch r.tokens[i] {
		case 36: // $ (Bulk strings)
			if parsedHeader {
				for isNumber(r.tokens[i+1]) {
					i++
				}
			}
		case 42: // * (Arrays)
			if !parsedHeader {
				var bs []byte
				for isNumber(r.tokens[i+1]) {
					bs = append(bs, r.tokens[i+1])
					i++
				}
				l, err := strconv.Atoi(string(bs))
				if err != nil {
					return err
				}
				argLength = l
				parsedHeader = true
			}
		case 43: // + (Simple strings)
			r.simple = true
		default:
			cur = append(cur, r.tokens[i])
		}
	}

	if len(r.args) != argLength {
		fmt.Printf("Args mimatch!")
	}

	return nil
}

func isNumber(b byte) bool {
	return b >= 48 && b <= 57
}
