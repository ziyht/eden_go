package ecache

import (
	"fmt"
	"time"

	"github.com/ziyht/eden_go/utils/ptr"
)

func toBytesKey(key any)([]byte, error){
	switch k := key.(type) {
	case string: return ptr.StringToBytes(k), nil
	case []byte: return k, nil
	}
	return nil, fmt.Errorf("invalid key type, only support string and []byte")
}

func stringsToBytesArr(ks []string)([][]byte) {
	out := make([][]byte, 0, len(ks))
	for _, k := range ks {
		out = append(out, ptr.StringToBytes(k))
	}
	return out
}

func toBytesArr(key any)([][]byte, error){
	switch k := key.(type) {
	case string  : return [][]byte{[]byte(k)}, nil
	case []byte  : return [][]byte{k}, nil
	case []string: return stringsToBytesArr(k), nil
	case [][]byte: return k, nil
	}
	return nil, fmt.Errorf("only support string, []string, []byte and [][]byte")
}

func toBytesKeyVal(key any, val any, ttl ...time.Duration) ([]byte, []byte, error) {
	var key_ []byte
	switch k := key.(type) {
		case string: key_ = ptr.StringToBytes(k)
		case []byte: key_ = k
		default: return nil, nil, fmt.Errorf("invalid key type, only support string and []byte")
	}
	
	switch v := val.(type) {
		case string: return key_, ptr.StringToBytes(v), nil
		case []byte: return key_, v, nil
	}
	return nil, nil, fmt.Errorf("invalid val type, only support string and []byte")
}
