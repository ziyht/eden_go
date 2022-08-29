package driver

import (
	"fmt"
	"net/url"
	"strings"
)

type Driver interface {
	Open(path string, params map[string][]string)(DB, error)
}

const VALID_FORMAT = "<driver>:<path>[?arg1=val1[&...]]"

var drivers = make(map[string]Driver)
var driverNames = make([]string, 0)

func Register(name string, driver Driver) error {
	if driver == nil {
		return fmt.Errorf("input driver is nil")
	}

	if _, ok := drivers[name]; ok {
		return fmt.Errorf("driver named '%s' already registered", name)
	} else {
		drivers[name] = driver
		driverNames = append(driverNames, name)
	}
	return nil
}

func OpenDsn(dsn string) (DB, error) {
	// get and find driver
	idx := strings.Index(dsn, ":")
	if idx < 0 {
		return nil, fmt.Errorf("invalid dsn(%s), the valid format is: %s", dsn, VALID_FORMAT)
	}
	driverName := dsn[:idx]
	driver := drivers[driverName]
	if driver == nil {
		return nil, fmt.Errorf("driver named '%s' can not be found, now support drivers are: %s", driverName, driverNames)
	}

	// parsing path
	pathStr := dsn[idx+1:]
	paramsStr := ""
	{
		idx := strings.Index(pathStr, "?")
		if idx == 0 {
			return nil, fmt.Errorf("can not find path in dsn(%s), the valid format is: %s", dsn, VALID_FORMAT)
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
			return nil, fmt.Errorf("parse params failed: %s, input dsn is '%s'", err, dsn)
		}
	}

	return driver.Open(pathStr, params)
}
