package driver

import (
	"fmt"
	"strings"
)

type Driver interface {
	Open(dsn string)(DB, error)
}

var drivers = make(map[string]Driver)
var driverNames = make([]string, 0)

func Register(name string, driver Driver) error {
	if _, ok := drivers[name]; ok {
		return fmt.Errorf("driver named '%s' already registered", name)
	} else {
		drivers[name] = driver
		driverNames = append(driverNames, name)
	}
	return nil
}

func OpenDsn(dsn string) (DB, error) {
	splits := strings.SplitN(dsn, ":", 2)
	if len(splits) != 2 {
		return nil, fmt.Errorf("invalid dsn(%s), the format should be: <dirver>://...", dsn)
	}

	d := drivers[splits[0]]
	if d == nil {
		return nil, fmt.Errorf("driver named '%s' can not be found, now support drivers are: %s", splits[0], driverNames)
	}
	
	return d.Open(splits[1])
}