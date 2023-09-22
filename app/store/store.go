package store

import (
	"errors"
	"fmt"
	"time"
)

var store map[string]*value

func init() {
	store = make(map[string]*value)
}

type value struct {
	Data string
	Exp  time.Time
}

func (v *value) SetExpiry(t time.Time) {
	v.exp = t
}

func Set(k, v string, px int) error {
	val := &value{Data: v}
	if px > 0 {
		val.SetExpiry(time.Now().Add(time.Millisecond * time.Duration(px)))
	}
	fmt.Printf("Val = %v\n", val)
	fmt.Printf("&Val = %v\n", &val)
	fmt.Printf("*Val = %v\n", *val)
	store[k] = val
	return nil
}

func Get(k string) (string, error) {
	val := *store[k]
	fmt.Printf("Val = %v\n", val)
	fmt.Printf("&Val = %v\n", &val)
	fmt.Println("CHECKING IS ZERO")
	if val.Exp.IsZero() {
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
	fmt.Printf("Checkign if %v is after %v", now, v.Exp)
	return now.After(v.Exp)
}
