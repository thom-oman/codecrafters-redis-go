package store

var store map[string]string

func init() {
	store = make(map[string]string)
}

func Set(k, v string) error {
	store[k] = v
	return nil
}

func Get(k string) (string, error) {
	return store[k], nil
}
