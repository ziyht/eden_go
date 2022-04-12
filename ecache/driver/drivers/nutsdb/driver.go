package pebble

import (

	"net/url"

	"github.com/ziyht/eden_go/ecache/driver"
)


type myDriver struct {
}

var driverName = "nutsdb"
var sampleDsn = "nutsdb:<dir>"
var insDriver = &myDriver{}

func (d *myDriver)Open(dsn string) (driver.DB, error) {

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	
	cfg := cfg{
		Dir: uri.Path,
	}

	return  newDB(&cfg)
}

func init() {
	driver.Register(driverName, insDriver)
}