package bencode

import "fmt"

func GetUint64(i interface{}) (uint64, error) {
	switch v := i.(type) {
	case int64:
		return uint64(v), nil
	case uint64:
		return v, nil
	default:
		return 0, fmt.Errorf("wrong type: %T", i)
	}
}
