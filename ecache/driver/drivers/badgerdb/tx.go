package badgerdb

import (
	"fmt"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type TX struct {
	txn *badger.Txn
}

func __genStoreKey(prefix []byte, setkey []byte)(storeKey []byte) {
	if len(prefix) == 0 {
		return setkey
	}

	out := make([]byte, 0, len(prefix) + len(setkey))
	out = append(out, prefix...)
	out = append(out, setkey...)
	return out
}

func (tx *TX)Set(prefix []byte, key []byte, val []byte, ttl ...time.Duration) error{
	e := badger.NewEntry(__genStoreKey(prefix, key), val)
	if len(ttl) > 0 && ttl[0] > 0 {
		e.WithTTL(ttl[0])
	}

	return tx.txn.SetEntry(e)
}

func (tx *TX)Get(prefix []byte, key []byte, del ...bool) ([]byte, uint64, error){
	e, err := tx.txn.Get(__genStoreKey(prefix, key))
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return nil, 0, err
		}
		return nil, 0, nil
	}

	val, err := e.ValueCopy(nil)
	if err != nil {
		return nil, 0, err
	}

	if len(del) > 0 && del[0] {
		if err = tx.txn.Delete(e.Key()); err != nil {
			return nil, 0, err
		}
	}

	return val, e.ExpiresAt(), nil
}

func (tx *TX)Del(prefix []byte, key []byte) (error){
	return tx.txn.Delete(__genStoreKey(prefix, key))
}

func (tx *TX)Iterate(prefix []byte, fn func(idx int, key []byte, val[]byte, expiresAt uint64)error) (error){
	it := tx.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	idx := -1
	prelen := len(prefix)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		idx += 1
		e := it.Item()
		key := e.KeyCopy(nil)
		val, err := e.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("ValurCopy failed: %s", err)
		}

		if err = fn(idx, key[prelen:], val, e.ExpiresAt()); err != nil {
			return err
		}
	}

	return nil
}