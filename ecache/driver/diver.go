package driver

import (
	"fmt"
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

func CheckDsn(dsn string) error {
	_, _, _, err := parsingDsn(dsn)
	return err
}

func OpenDsn(dsn string) (DB, error) {
	driver, pathStr, params, err := parsingDsn(dsn)
	if err != nil {
		return nil, err
	}

	return driver.Open(pathStr, params)
}

