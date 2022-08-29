package badgerdb

import (
	"github.com/ziyht/eden_go/ecache/driver"
)


type myDriver struct {
}

var driverName = "badger"
var insDriver  = &myDriver{}

func (d *myDriver)Open(path string, params map[string][]string) (driver.DB, error) {
	cfg := cfg{
		Dir     : path,
		InMemory: driver.GetBool(params, "memory") || driver.GetBool(params, "in-memory"),
	}

	return  newDB(&cfg)
}

func init() {
	driver.Register(driverName, insDriver)
}