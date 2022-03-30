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
var rpre = []byte{7, 127, 6} 
var bpre = []byte{7, 127, 7}
var _pre = []byte{6, 127}

// rpre + bucket
func bucketRecordKey(bucket string) []byte {
	out := []byte{};
	out = append(out, rpre...)
	out = append(out, []byte(bucket)...)
	return out
}

// bpre + bucket
func bucketWriteKeyPrefix(bucket string) []byte {
	out := []byte{};
	out = append(out, bpre...)
	out = append(out, []byte(bucket)...)
	return out
}

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

func (db *DB) Close () error{
	return db.db.Close()
}

func (db *DB)reloadBucketNames() error {
	return db.db.View(func(txn *badger.Txn) error {
    buckets := make(map[string][]byte)
		prefix := rpre[:]
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
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
		k := bucketRecordKey(bucket)
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
		key := []byte{}; key = append(key, _pre...); key = append(key, k...)
		elem := badger.NewEntry(key, v)
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
	var prefix []byte = _pre[:]  // for none bucket
	if bucket != nil {
		if err = db.registerBucketName(*bucket); err != nil{
			return fmt.Errorf("registerBucketName for '%s' failed: %s", *bucket, err)
		}
		prefix = bucketWriteKeyPrefix(*bucket)
	} 
	return db.db.Update(func(txn *badger.Txn) error {
		for i, k := range ks{
			var e *badger.Entry
			key := []byte{}; key = append(key, prefix...); key = append(key, k...)
			if lv == 1 {
				e = badger.NewEntry(key, vs[0])
			} else {
				e = badger.NewEntry(key, vs[i])
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

func(db *DB)getAll(bucket *string)(ks[][]byte, vs[][]byte, err error){
	var prefix []byte = _pre[:]   // for none bucket
	if bucket != nil {
		if len(*bucket) == 0 {
			return nil, nil, fmt.Errorf("empty bucket name")
		}
		prefix = bucketWriteKeyPrefix(*bucket)
	}
	prefixl := len(prefix)
	err = db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()[prefixl:]
			if err := item.Value(func(v []byte) error {
				ks = append(ks, key)
				vs = append(vs, v)
				return nil
			}); err != nil {
				return err
			}
		}
		
		return nil
	})
	return 
}

func(db *DB)Sets(ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	return db.sets(nil, ks, vs, durations...)
}

func (db *DB)SetIFs(items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	return db.db.Update(func(txn *badger.Txn) error {
		for idx, item := range items {
			k, v , d := fn(idx, item)
			key := []byte{}; key = append(key, _pre...); key = append(key, k...)
			elem := badger.NewEntry(key, v)
			if d > 0 { elem.WithTTL(d) }
			if err := txn.SetEntry(elem); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB)get(bucket *string, k []byte, del bool)(v []byte, err error){
	if bucket != nil {
		if len(*bucket) == 0 {
			return nil, fmt.Errorf("empty bucket name")
		}
	}

	fn := func(txn *badger.Txn) error {
		if bucket != nil {
			k = append(bucketWriteKeyPrefix(*bucket), k...)
		} else {
			k = append(_pre[:], k...)
		}
		elem, err := txn.Get(k)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			} else {
				return nil
			}
		}
		if err = elem.Value(func(val []byte) error {
			v = val
			return nil
		}); err != nil{
			return nil
		}
		
		if !del {
			return nil
		}

		return txn.Delete(k)
	}
	if !del {
		err = db.db.View(fn)
	} else {
		err = db.db.Update(fn)
	}
	
	return
}

func(db *DB)Get(k []byte, del bool) (v []byte, err error) {
	return db.get(nil, k, del)
}

func(db *DB)doForKeys(bucket *string, ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {
	var prefix []byte = _pre[:]
	if bucket != nil {
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}
		prefix = bucketWriteKeyPrefix(*bucket)
	}
	key := bytes.NewBuffer(nil)
	return db.db.View(func(txn *badger.Txn) error {
		for idx, k := range ks {
			key.Reset()
			key.Write(prefix)
			key.Write(k)
			elem, err := txn.Get(key.Bytes())
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
		}
		return nil
	})
}

func(db *DB)DoForKeys(ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {
	return db.doForKeys(nil, ks, fn)
}

func(db *DB)getVals(bucket *string, ks [][]byte, del bool) ([][]byte,  error) {
	if len(ks) == 0 {
		return nil, nil
	}
	out := make([][]byte, len(ks))
	fn := func(txn *badger.Txn) error {
		var prefix []byte = _pre[:]
		if bucket != nil {
			if len(*bucket) == 0 {
				return fmt.Errorf("empty bucket name")
			}
			prefix = bucketWriteKeyPrefix(*bucket)
		}

		var key *bytes.Buffer
		if !del {
			key = bytes.NewBuffer(nil)
		}
		for i, k := range ks {
			if key != nil {
				key.Reset()
				key.Write(prefix)
				key.Write(k)
				k = key.Bytes()
			} else {
				k_ := []byte{}; k_ = append(k_, prefix...);
				k = append(k_, k...)
			}
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

			if del {
				if err := txn.Delete(k); err != nil{
					return err
				}
			}
		}
		return nil
	}

	var err error
	if !del {
		err = db.db.View(fn)
	} else {
		err = db.db.Update(fn)
	}
	

	return out, err
}

func(db *DB)Gets(ks [][]byte, del bool) ([][]byte,  error) {
	return db.getVals(nil, ks, del)
}

func (db *DB)GetAll()([][]byte, [][]byte, error){
	return db.getAll(nil)
}

func (db *DB)Del(k []byte) error {
	return db.Dels([][]byte{k})
}

func (db *DB)Dels(ks [][]byte)(err error) {
	return db.db.Update(func(txn *badger.Txn) error {
		for _, k := range ks{
			key := []byte{}; key = append(key, _pre...); key = append(key, k...)
			if err := txn.Delete(key); err != nil{
				return err
			}
		}
		return nil
	})
}

func (db *DB)DelsAndGets(ks [][]byte)([][]byte, error) {
	return db.getVals(nil, ks, true)
}

func(db *DB)BSet(bucket string, k []byte, v []byte, duration ...time.Duration) (err error) {
	return db.sets(&bucket, [][]byte{k}, [][]byte{v}, duration...)
}

func(db *DB)BSets(bucket string, ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	return db.sets(&bucket, ks, vs, durations...)
}

func (db *DB)BSetIFs(bucket string, items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	if err := db.registerBucketName(bucket); err != nil {
		return err
	}
	prefix := bucketWriteKeyPrefix(bucket)
	var elems []*badger.Entry 
	return db.db.Update(func(txn *badger.Txn) error {
		for idx, item := range items {
			k, v , d := fn(idx, item)
			key := []byte{}; key = append(key, prefix...); key = append(key, k...)
			elem := badger.NewEntry(key, v)
			elems = append(elems, elem)
			if d > 0 { elem.WithTTL(d) }
			if err := txn.SetEntry(elem); err != nil {
				return err
			}
		}
		return nil
	})
}

func(db *DB)BGet(bucket string, k []byte, del bool) (v []byte, err error) {
	return db.get(&bucket, k, del)
}

func (db *DB)BGets(bucket string, ks [][]byte, del bool)(vs [][]byte, err error) {
	return db.getVals(&bucket, ks, del)
}

func (db *DB)BGetAll(bucket string)([][]byte, [][]byte, error){
	return db.getAll(&bucket)
}

func (db *DB)BDoForKeys(bucket string, ks [][]byte, fn func(idx int, key []byte, val []byte) error) error {
	return db.doForKeys(&bucket, ks, fn)
}

func(db *DB)doForAll(bucket *string, fn func(idx int, key []byte, val []byte) error) error {
	prefix := _pre[:]
	if bucket != nil{
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}
		if db.buckets[*bucket] == nil {
			return nil
		}
		prefix = bucketWriteKeyPrefix(*bucket)
	}
	prefixl := len(prefix)
	return db.db.View(func(txn *badger.Txn) error {
		idx := -1
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			idx += 1
			item := it.Item()
			key := item.Key()[prefixl:]
			if err := item.Value(func(v []byte) error {
				return fn(idx, key, v)
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func(db *DB)DoForAll(fn func(idx int, key []byte, val []byte) error) error {
	return db.doForAll(nil, fn)
}

func(db *DB)BDoForAll(bucket string, fn func(idx int, key []byte, val []byte) error) error {
	return db.doForAll(&bucket, fn)
}

func (db *DB)clear(bucket *string) error {
	if bucket != nil {
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}
		return db.db.DropPrefix(bucketWriteKeyPrefix(*bucket))
	}
	return db.db.DropPrefix(_pre[:])
}

func (db *DB)Clear() error {
	return db.clear(nil)
}

func (db *DB)BClear(bucket string) error {
	return db.clear(&bucket)
}

func (db *DB)BDel(bucket string, k []byte)(error) {
	if len(bucket) == 0 {
    return fmt.Errorf("empty bucket name")
	}
	key := append(bucketWriteKeyPrefix(bucket), k...)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (db *DB)BDels(bucket string, ks [][]byte)(error) {
	if len(bucket) == 0 {
    return fmt.Errorf("empty bucket name")
	}

	prefix := bucketWriteKeyPrefix(bucket)
	return db.db.Update(func(txn *badger.Txn) error {
		for _, k := range ks {
			key := []byte{}; key = append(key, prefix...); key = append(key, k...)
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func(db *DB)Truncate() (error) {
	return db.db.DropAll()
}