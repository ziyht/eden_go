package badgerdb

import (

	"net/url"

	"github.com/ziyht/eden_go/ecache/driver"
)


type myDriver struct {
}

var driverName = "badger"
var sampleDsn = "badger:<dir>?value_dir=<dir>&"
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