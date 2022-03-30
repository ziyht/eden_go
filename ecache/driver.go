package ecache

import (
	"fmt"
	"strings"
	"time"
)

type DB interface {
	Set(k []byte, v []byte, ttl ...time.Duration) error
	Sets(ks [][]byte, vs[][]byte, ttls ...time.Duration) error
	SetIFs(items []any, fn func(idx int, item any)(k[]byte, v[]byte, du time.Duration))error
	Get(k []byte, del bool) ([]byte, error)
	Gets(ks [][]byte, del bool)([][]byte, error)
	GetAll()([][]byte, [][]byte, error)
	Del(k []byte) error
	Dels(ks [][]byte) error
	DoForKeys(ks [][]byte, fn func(idx int, key []byte, val []byte) error) error
	DoForAll(fn func(idx int, key []byte, val []byte) error) error
	Clear() error     // cleat all keys, not include buckets
	
	Buckets()[]string                 // return all bucket names
	HaveBucket(bucket string)bool     // is the specific bucket exist
	BSet(bucket string, k []byte, v []byte, duration ...time.Duration) error
	BSets(bucket string, ks [][]byte, vs [][]byte, duration ...time.Duration) error
	BSetIFs(bucket string, items []any, fn func(idx int, item any)(k[]byte, v[]byte, du time.Duration)) error
	BGet(bucket string, k []byte, del bool) (v []byte, err error)
	BGets(bucket string, ks [][]byte, del bool)([][]byte, error)
	BGetAll(bucket string)([][]byte, [][]byte, error)
	BDel(bucket string, k []byte) error
	BDels(bucket string, k [][]byte) error
	BDoForKeys(bucket string, ks [][]byte, fn func(idx int, key []byte, val []byte) error) error
	BDoForAll(bucket string, fn func(idx int, key []byte, val []byte) error) error
	BClear(bucket string) error
	
	Truncate() error
	Close() error

	// UpdateVal()  // todo: only update val 
	// UpdateTTL()  // todo: only update TTL
}

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

func getDriver(name string)Driver{
	return drivers[name]
}

func openDsn(dsn string) (DB, error) {
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