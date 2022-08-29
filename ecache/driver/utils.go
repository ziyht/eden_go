package driver

func GetBool(params map[string][]string, key string)bool{
	vals := params[key]
	if len(vals) > 0 {
		switch vals[0]{
		case "false": return false
		case "False": return false
		case "0"    : return false
		case "true" : return true
		case "True" : return true
		case "1"    : return true
		}
	}

	return false
}
