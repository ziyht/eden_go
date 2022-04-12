package driver

import "time"

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