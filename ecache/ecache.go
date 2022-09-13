package ecache

import (
	_ "github.com/ziyht/eden_go/ecache/driver/drivers/badgerdb"
	_ "github.com/ziyht/eden_go/ecache/driver/drivers/nutsdb"
)

func NewDBCache(dsn string) (c *DBCache, err error) {
	db, err := newDB(dsn)
	if err != nil {
		return nil, err
	}

	c = newDBCache(db)

	return
}

func NewMemCache[T any](opts ...MemCacheOpts[T]) (*MemCache[T]) {
	if len(opts) > 0 {
		return newMemCache(opts[0])
	}

	return newMemCache(MemCacheOpts[T]{})
}

func NewTypedItemRegion[T Item](r* Region)(*ItemRegion[T]){
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

	return NewDBCache(cfg.Dsn)
}

func GetDBCache(name ...string)(c *DBCache, err error){
	if len(name) == 0 {
		return getDfDBcache()
	}

	return getDBCache(name[0])
}
