package ecache

import (
	"fmt"
	"strings"

	_ "github.com/ziyht/eden_go/ecache/driver/drivers/badgerdb"
	_ "github.com/ziyht/eden_go/ecache/driver/drivers/nutsdb"
)

const (
	BADGER = "badger"
	NUTSDB = "nutsdb"
)

type DBCacheOpts struct {
	Dsn     string   // if set, the other proporty will take no effect, format: <dirvername>:<dir>[?<arg1=val1>[&<arg2=val2>]...]
	Driver  string   // if not set will using nutsdb in default
	Dir     string   
  Params  map[string][]string
}

func NewDBCache(opts DBCacheOpts) (c *DBCache, err error) {
	db, err := newDB(&opts)
	if err != nil {
		return nil, err
	}

	c = newDBCache(db)

	return
}

func NewMemCache[K Key, V any](opts ...MemCacheOpts[V]) (*MemCache[K, V]) {
	if len(opts) > 0 {
		return newMemCache[K, V](opts[0])
	}

	return newMemCache[K, V](MemCacheOpts[V]{})
}

func GenDsn(driver, dir string, params ...map[string][]string) string {
	dsn := NUTSDB
	if len(driver) > 0 {
		dsn = driver
	}
	dsn = dsn + ":" + dir
	var args []string
	for _, p := range params {
		for k, vs := range p {
			for _, v := range vs {
				args = append(args, fmt.Sprintf("%s=%s", k, v))
			}
		}
	}

	if len(args) > 0 {
		dsn += "?"
		dsn += strings.Join(args, "&")
	}

	return dsn
}



func MakeTypedItemRegion[T Item](r* Region)(*ItemRegion[T]){
	return newItemRegion[T](r.db, r.meta.keys)
}

func NewTypedItemRegion[T Item](r* Region, keys ...string)(*ItemRegion[T]){
	if len(keys) == 0 {
		return MakeTypedItemRegion[T](r)
	}

	r = r.SubRegion(keys...)
	return newItemRegion[T](r.db, r.meta.keys)
}

// InitFromConfigFile will init dbcache from a config file, support multi file types like yaml, yml, json, toml...
// 
// the format should like follows:
// ecache:
//   c1:
//     dsn: badger:./cache/ecache/c1
//
// after init, you can get dbcache by call GetDBCache("c1")
func InitFromConfigFile(path string)error{
	return initFromFile(path)
}

func GetDBCacheFromFile(path string, key string)(c *DBCache, err error){
	cfg, err := cfgsFromFileKey(path, key)
	if err != nil {
		return nil, err
	}

	return NewDBCache(DBCacheOpts{Dsn: cfg.Dsn} )
}

func GetDBCache(name ...string)(c *DBCache, err error){
	if len(name) == 0 {
		return getDfDBcache()
	}

	return getDBCache(name[0])
}
