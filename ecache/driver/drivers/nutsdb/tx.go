package nutsdb

import (
	"time"

	"github.com/xujiajun/nutsdb"
)

type TX struct {
	txn *nutsdb.Tx
}

func __validTTL(ttl ...time.Duration) uint32 {
	if len(ttl) > 0{
		secs := uint32(ttl[0].Seconds())
		if secs == 0 {
			if ttl[0] > 0 {
				secs = 1
			}
		}
		return secs
	}

	return 0
}

func (tx *TX)Set(prefix []byte, key []byte, val []byte, ttl ...time.Duration) error{
	return tx.txn.Put(string(prefix), key, val, __validTTL(ttl...))
}

func (tx *TX)Get(prefix []byte, key []byte, del ...bool) ([]byte, uint64, error){
	e, err := tx.txn.Get(string(prefix), key)
	if err != nil {
		if err != nutsdb.ErrKeyNotFound{
			return nil, 0, err
		}
	}

	if len(del) > 0 && del[0] {
		err = tx.txn.Delete(string(prefix), key)
		if err != nil {
			return nil, 0, err
		}
	}

	return e.Value, e.Meta.Timestamp + uint64(e.Meta.TTL), nil
}

func (tx *TX)Del(prefix []byte, key []byte) (error){
	return tx.txn.Delete(string(prefix), key)
}

func (tx *TX)Iterate(prefix []byte, fn func(idx int, key []byte, val []byte, expiresAt uint64)error) (error){
	es, err := tx.txn.GetAll(string(prefix))
	if err != nil {
		if err != nutsdb.ErrBucketEmpty{
			return err
		}
	}
	
	for i, e := range es {
		err = fn(i, e.Key, e.Value, e.Meta.Timestamp + uint64(e.Meta.TTL))
		if err != nil {
			return err
		}
	}

	return nil
}