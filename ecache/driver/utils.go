package driver

import (
	"fmt"
	"net/url"
	"strings"
)

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

func parsingDsn(dsn string) (Driver, string, url.Values, error){
	// get and find driver
	idx := strings.Index(dsn, ":")
	if idx < 0 {
		return nil, "", nil, fmt.Errorf("invalid dsn(%s), the valid format is: %s", dsn, VALID_FORMAT)
	}
	driverName := dsn[:idx]
	driver := drivers[driverName]
	if driver == nil {
		return nil, "", nil, fmt.Errorf("driver named '%s' can not be found, now support drivers are: %s", driverName, driverNames)
	}

	// parsing path
	pathStr := dsn[idx+1:]
	paramsStr := ""
	{
		idx := strings.Index(pathStr, "?")
		if idx == 0 {
			return nil, "", nil, fmt.Errorf("can not find path in dsn(%s), the valid format is: %s", dsn, VALID_FORMAT)
		} else if idx > 0{
			paramsStr = pathStr[idx+1:]
			pathStr   = pathStr[:idx]
		}
	}

	// parsing args
	var params url.Values
	var err error
	if paramsStr != "" {
		params, err = url.ParseQuery(paramsStr)
		if err != nil {
			return nil, "", nil, fmt.Errorf("parse params failed: %s, input dsn is '%s'", err, dsn)
		}
	}

	return driver, pathStr, params, nil
}