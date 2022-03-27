package pebble

import (
	"bytes"
	"fmt"
	"time"
	"golang.org/x/exp/maps"

	badger "github.com/dgraph-io/badger/v3"
)

type cfg struct {
	Dir  string
}

type DB struct {
	db      *badger.DB
	buckets map[string][]byte
}

type TX struct {
	txn *badger.Txn
}

// magic val for set bucket
var bmagic  = []byte{1, 127, 1}
var bprefix = []byte{1, 127, 1, 0}

func newDB(cfg *cfg) (*DB, error){

	opt := badger.DefaultOptions(cfg.Dir)

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	out := &DB{db: db}
	out.buckets = make(map[string][]byte)
	out.reloadBucketNames()

	return out, nil
}

func bucketKey(bucket string) []byte {
	key := bprefix[:]
	key = append(key, []byte(bucket)...)
	return key
}

func (db *DB) Close () error{
	return db.db.Close()
}

func (db *DB)reloadBucketNames() error {
	return db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := bprefix
		buckets := make(map[string][]byte)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				buckets[string(v)] = k
				return nil
			})
			if err != nil {
				return err
			}
		}
		db.buckets = buckets
		return nil
	})
}

func (db *DB)registerBucketName(bucket string) error {
	if len(bucket) == 0 {
		return fmt.Errorf("bucket name should not be a empty string")
	}

	// bucket already exist
	if db.buckets[bucket] != nil{
		return nil
	}

	return db.db.Update(func(txn *badger.Txn) error {
		k := bucketKey(bucket)
		err := txn.SetEntry(badger.NewEntry(k, []byte(bucket)))
		if err != nil {
			return err
		}
		db.buckets[bucket] = k
		return nil
	})
}

func (db *DB)Buckets()[]string{
	return maps.Keys(db.buckets)
}

func (db *DB)HaveBucket(bucket string)bool{
	_, ok := db.buckets[bucket]; return ok
}

func (db *DB)NewTransaction()*TX{
	return &TX{txn: db.db.NewTransaction(false)}
}

func(db *DB)Set(k []byte, v []byte, duration ...time.Duration) (err error) {
	return db.db.Update(func(txn *badger.Txn) error {
		elem := badger.NewEntry(k, v)
		if len(duration) > 0 {
			elem.WithTTL(duration[0])
		}

		return txn.SetEntry(elem)
	})
}

func(db *DB)sets(bucket *string, ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	lk := len(ks); lv := len(vs); ld := len(durations);
	if lv != lk {
		if lv != 1 {
			return fmt.Errorf("the len for key(%d) and values(%d) are not match", lk, lv)
		}
	}
	if ld != 0 && ld != lv {
		if ld != 1 {
			return fmt.Errorf("the len for key(%d) and durations(%d) are not match", lk, ld)
		}
	}
	prefix := []byte{}
	if bucket != nil {
		if err = db.registerBucketName(*bucket); err != nil{
			return fmt.Errorf("registerBucketName for '%s' failed: %s", *bucket, err)
		}
		prefix = append(prefix, []byte(*bucket)...)
		prefix = append(prefix, bmagic...)
	}
	key := bytes.NewBuffer(nil)
	return db.db.Update(func(txn *badger.Txn) error {
		for i := 0; i < lk; i++ {
			var e *badger.Entry
			key.Reset()
			key.Write(prefix)
			key.Write(ks[i])
			if lv == 1 {
				e = badger.NewEntry(key.Bytes(), vs[0])
			} else {
				e = badger.NewEntry(key.Bytes(), vs[i])
			}
			
			if ld == 1 {
				e.WithTTL(durations[0])
			} else if ld > 1 {
				e.WithTTL(durations[i])
			}

			err := txn.SetEntry(e)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func(db *DB)Sets(ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	return db.sets(nil, ks, vs, durations...)
}

func(db *DB)BSet(bucket string, k []byte, v []byte, duration ...time.Duration) (err error) {
	return db.sets(&bucket, [][]byte{k}, [][]byte{v}, duration...)
}

func(db *DB)BSets(bucket string, ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	return db.sets(&bucket, ks, vs, durations...)
}

func (db *DB)SetIFs(items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	return db.db.Update(func(txn *badger.Txn) error {
		for idx, item := range items {
			k, v , d := fn(idx, item)
			elem := badger.NewEntry(k, v)
			if d > 0 { elem.WithTTL(d) }
			if err := txn.SetEntry(elem); err != nil {
				return err
			}
		}
		return nil
	})
}

func(db *DB)Get(k []byte) (v []byte, err error) {
	err = db.db.View(func(txn *badger.Txn) error {
		elem, err := txn.Get(k)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		elem.Value(func(val []byte) error {
			v = val
			return nil
		})
		
		return nil
	})

	return 
}

func(db *DB)DoForKeys(ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {

	return db.db.View(func(txn *badger.Txn) error {
		for idx, k := range ks {
			elem, err := txn.Get(k)
			if err == badger.ErrKeyNotFound{
				if err = fn(idx, k, nil); err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
			if err = elem.Value(func(val []byte) error {
				return fn(idx, k, val)
			}); err != nil {
				return err
			}
			
			return nil
		}
		return nil
	})
}

func(db *DB)getVals(ks [][]byte, delete bool) ([][]byte,  error) {
	if len(ks) == 0 {
		return nil, nil
	}
	out := make([][]byte, len(ks))
	err := db.db.View(func(txn *badger.Txn) error {
		for i, k := range ks {
			elem, err := txn.Get(k)
			if err == badger.ErrKeyNotFound {
				out[i] = nil
				continue
			}
			if err != nil {
				out = nil
				return err
			}

			if err = elem.Value(func(val []byte) error {
				out[i] = val
				return nil
			}); err != nil {
				return err
			}

			if delete {
				if err := txn.Delete(k); err != nil{
					return err
				}
			}
		}
		return nil
	})

	return out, err
}

func(db *DB)Gets(ks [][]byte) ([][]byte,  error) {
	return db.getVals(ks, false)
}

func (db *DB)Del(k []byte) error {
	return db.Dels([][]byte{k})
}

func (db *DB)Dels(ks [][]byte)(err error) {
	return db.db.Update(func(txn *badger.Txn) error {
		for _, k := range ks{
			if err := txn.Delete(k); err != nil{
				return err
			}
		}
		return nil
	})
}

func (db *DB)DelsAndGets(ks [][]byte)([][]byte, error) {
	return db.getVals(ks, true)
}

func(db *DB)BGet(bucket string, k []byte) (v []byte, err error) {
	key := append([]byte(bucket), bmagic...)
	key = append(key, k...)
	return db.Get(key)
}

func (db *DB)BDel(k []byte)([]byte, error) {
	return nil, fmt.Errorf("todo")
}

func(db *DB)Clear() (error) {
	return db.db.DropAll()
}