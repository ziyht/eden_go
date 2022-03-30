package pebble

import (

	"net/url"

	"github.com/ziyht/eden_go/ecache"
)


type driver struct {
}

var driverName = "badger"
var sampleDsn = "badger:<dir>?value_dir=<dir>&"
var insDriver = &driver{}

func (d *driver)Open(dsn string) (ecache.DB, error) {

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
	ecache.Register(driverName, insDriver)
}
