package parser

import (
	"testing"
)

func TestAddToken(t *testing.T) {
	tests := []struct {
		name         string
		tokens       []byte
		expectedSize int
	}{
		{
			"normal",
			[]byte("abc"),
			3,
		},
		{
			"ignores 0 bytes",
			append([]byte("abc"), 0),
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest()
			if req.TokenCount() != 0 {
				t.Fatalf("Request should start with size 0, got %v", req.TokenCount())
			}

			for _, b := range tt.tokens {
				req.AddToken(b)
			}
			if req.TokenCount() != tt.expectedSize {
				t.Fatalf("Expected size of %v but got %v", tt.expectedSize, req.TokenCount())
			}

		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		tokens []byte
		args   []string
		simple bool
		err    error
	}{
		{
			"valid ping command",
			[]byte("*1\r\n$4\r\nping\r\n"),
			[]string{"ping"},
			false,
			nil,
		},
		{
			"valid echo command",
			[]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"),
			[]string{"ECHO", "hey"},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest()
			req.AddTokens(tt.tokens)
			err := req.Parse()
			// fmt.Printf("Parsed args: %v", req.Args())

			if err != tt.err {
				if tt.err != nil {
					t.Fatalf("Expected to get error %v but instead got %\n", tt.err, err)
				} else {
					t.Fatalf("Expected to get not error but instead got %v\n", err)
				}
			}

			if len(req.Args()) != len(tt.args) {
				t.Fatalf("Expected to receive %v arguments but instead got %v\n", len(tt.args), len(req.args))
			}

			for i := range tt.args {
				args := req.Args()
				if tt.args[i] != args[i] {
					t.Fatalf("Expected element %v to match %v, received %v\n", i, tt.args[i], args[i])
				}
			}
		})
	}
}
