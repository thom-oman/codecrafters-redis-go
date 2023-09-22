package store

import (
	"errors"
	"fmt"
	"time"
)

var store map[string]value

func init() {
	store = make(map[string]value)
}

type value struct {
	Data string
	exp  time.Time
}

func Set(k, v string, px int) error {
	val := value{Data: v}
	if px > 0 {
		val.exp = time.Now().Add(time.Millisecond * time.Duration(px))
	}
	store[k] = val
	return nil
}

func Get(k string) (string, error) {
	val := store[k]
	if val.exp.IsZero() {
		return val.Data, nil
	}

	fmt.Println("CHECKING EXPIRY")
	if expired(val) {
		return "", errors.New("Key has expired")
	}
	fmt.Println("CHECKED EXPIRY")

	return val.Data, nil
}

func expired(v value) bool {
	now := time.Now()
	fmt.Printf("Checkign if %v is after %v", now, v.exp)
	return now.After(v.exp)
}
