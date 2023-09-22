package store

var store map[string]string

func init() {

}

func Set(k, v string) error {
	store[k] = v
	return nil
}

func Get(k string) (string, error) {
	return store[k], nil
}
