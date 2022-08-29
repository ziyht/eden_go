package ecache

import (
	_ "github.com/ziyht/eden_go/ecache/driver/drivers/badgerdb"
	//_ "github.com/ziyht/eden_go/ecache/driver/drivers/nutsdb"
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