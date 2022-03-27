package ecache

import (
	"fmt"
	"time"
)

type ECache struct {
	dsn  string
	db   DB
}

func GetCache(name string) *ECache {
	return nil
}

func GetCacheFromDsn(dsn string )(*ECache, error){

	db, err := openDsn(dsn)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("invalid returned db(nil) checked from current driver in dsn(%s)", dsn)
	}

	return &ECache{dsn: dsn, db: db}, nil
}

func (c *ECache)Buckets(bucket string) []string {
	return c.db.Buckets()
}

func (c *ECache)HaveBucket(bucket string) bool {
	return c.db.HaveBucket(bucket)
}

// set <k, v> to cache with TTL
// if TTL not set, it will never expire
func (c *ECache)Set(k []byte, v []byte, ttl... time.Duration) error {
	return c.db.Set(k, v, ttl...)
}

func (c *ECache)Sets(ks [][]byte, vs [][]byte, ttls... time.Duration) error {
	return c.db.Sets(ks, vs, ttls...)
}

// set <k, v> to cache with TTL
// if TTL not set, it will never expire
func (c *ECache)BSet(bucket string, k []byte, v []byte, ttl... time.Duration) error {
	return c.db.BSet(bucket, k, v, ttl...)
}

func (c *ECache)BGet(bucket string, k []byte)([]byte, error) {
	return c.db.BGet(bucket, k)
}

func (c *ECache)MSetIFs(items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	return c.db.SetIFs(items, fn)
}

// get value by specified k
// note: it should not return error if k not exist in cache
func (c *ECache)Get(k []byte) ([]byte, error) {
	return c.db.Get(k)
}

// get values by specified keys
// the len of returned value will be the same as keys, the value will be nil in the same position as key in keys if key not exist in the cache
func (c *ECache)MGet(ks [][]byte) ([][]byte, error) {
	return c.db.Gets(ks)
}

func (c ECache)Clear() error {
	return c.db.Clear()
}

func (c *ECache)Close() error {
	return c.db.Close()
}

