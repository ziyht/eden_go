package driver

import "time"

type DB interface {	
	TX(tx interface{}) TX
	Update(func(tx TX)error) error
	View(func(tx TX)error) error

	DropPrefix(prefix []byte) error
	Truncate() error
	Close() error
}

type TX interface{
	/* 
	    note: real key will be [prefix + key] and prefix can be nil
	*/


  // set key val to db, ttl should only be set when len(ttl)>0 && ttl[0]>0
	Set(prefix []byte, key []byte, val []byte, ttl ...time.Duration) error

	// find the data and expiresAt(unit timestamp) for the key from db, and it will be deleted when len(del)>0 && del[0]==true
	// return (nil, 0, nil) if the key is not exist
	Get(prefix []byte, key []byte, del... bool)([]byte, uint64, error)

	// delete the key from db, return nil if the key is not exist
	Del(prefix []byte, key []byte) error

	// iterate all the keys have the same prefix, the the feed to fn is not been cut off prefix, you can do this operation by you self
	// the key passed in fn has been trimed out the prefix
	Iterate(prefix []byte, fn func(idx int, key []byte, val []byte, expiredAt uint64)error) error
}