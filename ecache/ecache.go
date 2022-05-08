package ecache

import (
	"fmt"
	"sync"
	"time"

	"github.com/ziyht/eden_go/ecache/driver"

	_ "github.com/ziyht/eden_go/ecache/driver/drivers/badgerdb"
	_ "github.com/ziyht/eden_go/ecache/driver/drivers/nutsdb"
)



type ECache struct {
	dsn  string
	db   driver.DB
	l1s  map[string]*ECacheL1
	mu   sync.Mutex
}

func GetCache(name string) *ECache {
	return nil
}

func GetCacheFromDsn(dsn string )(*ECache, error){

	db, err := driver.OpenDsn(dsn)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("invalid returned db(nil) checked from current driver in dsn(%s)", dsn)
	}

	return &ECache{dsn: dsn, db: db}, nil
}

// set <k, v> to cache with TTL
// if TTL not set, it will never expire
func (c *ECache)Set(k []byte, v []byte, ttl... time.Duration) error {
	return c.db.Set(k, v, ttl...)
}

func (c *ECache)Sets(ks [][]byte, vs [][]byte, ttls... time.Duration) error {
	return c.db.Sets(ks, vs, ttls...)
}

func (c *ECache)SetIFs(items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	return c.db.SetIFs(items, fn)
}

// get value by specified k
// note: it should not return error if k not exist in cache
func (c *ECache)Get(k []byte) ([]byte, error) {
	return c.db.Get(k, false)
}

// get values by specified keys
// the len of returned value will be the same as keys, the value will be nil in the same position as key in keys if key not exist in the cache
func (c *ECache)Gets(ks ...[]byte) ([][]byte, error) {
	return c.db.Gets(ks, false)
}

func (c *ECache)Del(k []byte) error {
	return c.db.Del(k)
}

func (c *ECache)Dels(ks ...[]byte) error {
	return c.db.Dels(ks)
}

func (c *ECache)GetAndDel(k []byte)([]byte, error){
	return c.db.Get(k, true)
}

func (c *ECache)GetsAndDel(ks ...[]byte)([][]byte, error){
	return c.db.Gets(ks, true)
}

func (c *ECache)DoForKeys(ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {
	return c.db.DoForKeys(ks, fn)
}

func (c *ECache)DoForAll(fn func(idx int, key []byte, val []byte) error) error {
	return c.db.DoForAll(fn)
}

func (c *ECache)Clear() error {
	return c.db.Clear()
}

func (c *ECache)Buckets(bucket string) []string {
	return c.db.Buckets()
}

func (c *ECache)HaveBucket(bucket string) bool {
	return c.db.HaveBucket(bucket)
}

// set <k, v> to cache with TTL
// if TTL not set, it will never expire
func (c *ECache)BSet(bucket string, k []byte, v []byte, ttl... time.Duration) error {
	return c.db.BSet(bucket, k, v, ttl...)
}

func (c *ECache)BSets(bucket string, ks [][]byte, vs [][]byte, ttls... time.Duration) error {
	return c.db.BSets(bucket, ks, vs, ttls...)
}

func (c *ECache)BSetIFs(bucket string, items []any, fn func(idx int, item any)(k[]byte, v[]byte, du time.Duration)) error {
	return c.db.BSetIFs(bucket, items, fn)
}

func (c *ECache)BGet(bucket string, k []byte)([]byte, error) {
	return c.db.BGet(bucket, k, false)
}

func (c *ECache)BGets(bucket string, ks ...[]byte)([][]byte, error) {
	return c.db.BGets(bucket, ks, false)
}

func (c *ECache)BDel(bucket string, k []byte) error {
	return c.db.BDel(bucket, k)
}

func (c *ECache)BDels(bucket string, ks ...[]byte) error {
	return c.db.BDels(bucket, ks)
}

func (c *ECache)BDoForKeys(bucket string, ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {
	return c.db.BDoForKeys(bucket, ks, fn)
}

func (c *ECache)BDoForAll(bucket string, fn func(idx int, key []byte, val []byte) error) error {
	return c.db.BDoForAll(bucket, fn)
}

func (c *ECache)BClear(bucket string) error {
	return c.db.BClear(bucket)
}

func (c *ECache)Truncate() error {
	return c.db.Truncate()
}

func (c *ECache)Close() error {
	return c.db.Close()
}

type Item interface {
	New() Item                // using this to make new Item in Get Methods
	Marshal()([]byte, error)  // do marshal things
	Unmarshal([]byte) error   // ummarshal data to self
	TTL() time.Duration       // return the TTL of current item, has lower priority than params
}

func (c *ECache)SetItem(k []byte, i Item, ttl ...time.Duration) error{
	v, err := i.Marshal()
	if err != nil { return err }

	if len(ttl) == 0{
		ittl := i.TTL()
		if ittl > 0 { ttl = []time.Duration{ittl} }
	}

	return c.Set(k, v, ttl...)
}

func (c *ECache)GetItem(k []byte, i Item)(out Item, err error){
	v, err := c.Get(k)
	if err != nil { return nil, err }
	if v   == nil { return nil, err }

	out = i.New()
	err = out.Unmarshal(v); if err != nil {out = nil}
	return
}

// BSetItem
// the ttl in params has higher priority than item.TTL()
func (c *ECache)BSetItem(bucket string, k []byte, i Item, ttl ...time.Duration) error{
	v, err := i.Marshal()
	if err != nil { return err }

	if len(ttl) == 0{
		ittl := i.TTL()
		if ittl > 0 { ttl = []time.Duration{ittl} }
	}

	return c.BSet(bucket, k, v, ttl...)
}

func (c *ECache)BGetItem(bucket string, k []byte, i Item)(out Item, err error) {
	v, err := c.BGet(bucket, k)
	if err != nil { return nil, err }
	if v   == nil { return nil, err }

	out = i.New()
	err = out.Unmarshal(v); if err != nil {out = nil}
	return
}