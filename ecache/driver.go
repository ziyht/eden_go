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
	
	// get the val by key, return nil, nil if key not exist
	Get(k []byte) ([]byte, error)
	Gets(ks [][]byte)([][]byte, error)
	DoForKeys(ks [][]byte, fn func(idx int, key []byte, val []byte) error) error

	// dels the vals of keys in cache
	Dels(ks [][]byte) (error)
	// dels the vals of keys in cache and return the values for which keys in the cache, the returned values will have the same len of keys, if some key not found, the same pos in val will set to nil            
	DelsAndGets(ks [][]byte)([][]byte, error)

	// UpdateVal()  // todo: only update val 
	// UpdateTTL()  // todo: only update TTL

	Buckets()[]string                 // return all bucket names
	HaveBucket(bucket string)bool     // is the specific bucket exist
	BSet(bucket string, k []byte, v []byte, duration ...time.Duration) error
	BGet(bucket string, k []byte) (v []byte, err error)
	BDel(k []byte)([]byte, error)
	//BDoForAll(fn func(idx int, key []byte, val []byte) error) error
	//ClearBucket() error

	//DelBucket()error
	

	Clear() error // truncate all values in cache
	Close() error
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