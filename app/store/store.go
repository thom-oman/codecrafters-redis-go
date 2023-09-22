package store

import (
	"errors"
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
	v.Exp = t
}

func Set(k, v string, px int) error {
	val := &value{Data: v}
	if px > 0 {
		val.SetExpiry(time.Now().Add(time.Millisecond * time.Duration(px)))
	}
	store[k] = val
	return nil
}

func Get(k string) (value string, err error) {
	val := *store[k]
	if val.Exp.IsZero() {
		value = val.Data
		return
	}

	if expired(val) {
		err = errors.New("Key has expired")
		return
	}

	return val.Data, nil
}

func expired(v value) bool {
	now := time.Now()
	return now.After(v.Exp)
}
