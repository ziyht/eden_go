package elog

import (
	"fmt"
	"reflect"
)


func getMultiBoolFromObj(obj interface{}) (b1 bool, b2 bool, err error) {
	switch obj := obj.(type) {
  case bool: return obj, false, nil
	case []interface{}: 
		if len(obj) == 0 { return false, false, fmt.Errorf("val not set")
	  } else {
			obj = append(obj, obj[0])
			var v1, v2 bool; var ok bool
			if v1, ok = obj[0].(bool); !ok {
				return false, false, fmt.Errorf("invalid type of [0]: %s", reflect.TypeOf(obj[0]))
			}
			if v2, ok = obj[1].(bool); !ok {
				return false, false, fmt.Errorf("invalid type of [1]: %s", reflect.TypeOf(obj[0]))
			}
			return v1, v2, nil
		}
  }  

	return false, false, fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}

func getMultiStringFromObj(obj interface{}) (out []string, err error) {
	switch obj := obj.(type) {
  case string: return []string{obj}, nil
	case []string:
		return obj, nil
	case []*string:
		for _, s := range obj { out = append(out, *s) }
		return 
	case []interface{}: 
		for i, o := range obj { 
			if s, ok := o.(string); ok {
				out = append(out, s)
			} else {
				return nil, fmt.Errorf("invalid type of [%d]: %s", i, reflect.TypeOf(obj[0]))
			}
		}
  }  
	return 
}

func getIntFromObj(obj interface{}) (val int, err error) {

	if v, ok := obj.(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}