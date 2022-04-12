package pebble

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"golang.org/x/exp/maps"

	nutsdb "github.com/xujiajun/nutsdb"
)

type cfg struct {
	Dir  string
}

type DB struct {
	db      *nutsdb.DB
	opts    *nutsdb.Options
	buckets map[string][]byte
}

type TX struct {
	txn *nutsdb.Tx
}

// magic val for set bucket
var rpre = []byte{7, 127, 6} 
var bpre = []byte{7, 127, 7}
var _pre = []byte{6, 127}

var empty_bucket_name  = ""
var record_bucket_name = "__buckets__"

// rpre + bucket
func bucketRecordKey(bucket string) []byte {
	out := []byte{};
	out = append(out, rpre...)
	out = append(out, []byte(bucket)...)
	return out
}

func newDB(cfg *cfg) (*DB, error){

	opts := nutsdb.DefaultOptions
	opts.Dir        = cfg.Dir
	opts.RWMode     = nutsdb.MMap
	opts.SyncEnable = false
	opts.SegmentSize= 64 * 1024 * 1024

	db, err := nutsdb.Open(opts)
	if err != nil {
		return nil, err
	}

	out := &DB{db: db}
	out.opts    = &opts
	out.buckets = make(map[string][]byte)
	out.reloadBucketNames()

	return out, nil
}

func (db *DB) Close () error{
	return db.db.Close()
}

func (db *DB)reloadBucketNames() error {
	return db.db.View(func(txn *nutsdb.Tx) error {
    buckets := make(map[string][]byte)
		entries, err := txn.GetAll(record_bucket_name)
		if err != nil{
			return nil
		}
		for _, entry := range entries{
			buckets[string(entry.Key)] = entry.Value
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

	return db.db.Update(func(txn *nutsdb.Tx) error {
		err := txn.Put(record_bucket_name, []byte(bucket), []byte(bucket), 0)
		if err != nil {
			return err
		}
		db.buckets[bucket] = []byte(bucket)
		return nil
	})
}

func (db *DB)Buckets()[]string{
	return maps.Keys(db.buckets)
}

func (db *DB)HaveBucket(bucket string)bool{
	_, ok := db.buckets[bucket]; return ok
}

func (db *DB)validTTL(duration ...time.Duration) uint32 {
	var ttl uint32
	if len(duration) > 0 {
		ttl = uint32(duration[0].Seconds())
		if ttl == 0 && duration[0] > 0 {
			ttl = 1
		}
	}
	return ttl
}

func(db *DB)Set(k []byte, v []byte, duration ...time.Duration) (err error) {
	return db.db.Update(func(txn *nutsdb.Tx) error {
		return txn.Put(empty_bucket_name, k, v, db.validTTL(duration...))
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
	if bucket != nil {
		if err = db.registerBucketName(*bucket); err != nil{
			return fmt.Errorf("registerBucketName for '%s' failed: %s", *bucket, err)
		}
	} else {
		bucket = &empty_bucket_name
	}

	for i := 0; i < lk;  {
		if err = db.db.Update(func(txn *nutsdb.Tx) error {
			cnt := 0
			for ; i < lk; i++ {
				if cnt == 1000 {
					return nil
				}

				ttl := uint32(0)
				if ld == 1 { ttl = db.validTTL(durations[0]) } else 
				if ld >  1 { ttl = db.validTTL(durations[i]) }

				if lv == 1 { err = txn.Put(*bucket, ks[i], vs[0], ttl) 
				} else     { err = txn.Put(*bucket, ks[i], vs[i], ttl)}

				if err != nil {
					fmt.Printf("%d\n\n", i)
					return err
				}
				cnt += 1
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func(db *DB)getAll(bucket *string)(ks[][]byte, vs[][]byte, err error){
	if bucket != nil {
		if len(*bucket) == 0 {
			return nil, nil, fmt.Errorf("empty bucket name")
		}
	} else {
		bucket = &empty_bucket_name
	}

	err = db.db.View(func(txn *nutsdb.Tx) error {

		entires, err := txn.GetAll(*bucket)
		if err != nil {
			return err
		}

		for _, ent := range entires {
			ks = append(ks, ent.Key)
			vs = append(vs, ent.Value)
		}
		
		return nil
	})
	return 
}

func(db *DB)Sets(ks [][]byte, vs [][]byte, durations ...time.Duration) (err error) {
	return db.sets(nil, ks, vs, durations...)
}

func (db *DB)setIFs(bucket *string, items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	if bucket != nil{
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}
		if err := db.registerBucketName(*bucket); err != nil {
			return err
		}
	} else {
		bucket = &empty_bucket_name
	}
	bucket_ := *bucket
	len := len(items)
	for i := 0; i < len; {
		if err := db.db.Update(func(txn *nutsdb.Tx) error {
			cnt := 0
			for ; i < len; i++ {
				if cnt == 1000 {
					return nil
				}
				k, v , d := fn(i, items[i])
				if err := txn.Put(bucket_, k, v, db.validTTL(d)); err != nil {
					return err
				}
				cnt += 1
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB)SetIFs(items []interface{}, fn func(int, interface{})(k[]byte, v[]byte, du time.Duration))error{
	return db.setIFs(nil, items, fn)
}

func (db *DB)get(bucket *string, k []byte, del bool)(v []byte, err error){
	if bucket != nil {
		if len(*bucket) == 0 {
			return nil, fmt.Errorf("empty bucket name")
		}
	} else {
		bucket = &empty_bucket_name
	}

	fn := func(txn *nutsdb.Tx) error {

		elem, err := txn.Get(*bucket, k)
		if err != nil {
			if err != nutsdb.ErrKeyNotFound {
				return err
			} else {
				return nil
			}
		}
		v = elem.Value
		
		if !del {
			return nil
		}

		return txn.Delete(*bucket, k)
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
	if bucket != nil {
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}		
	} else {
		bucket = &empty_bucket_name
	}
	bucket_ := *bucket

	return db.db.View(func(txn *nutsdb.Tx) error {
		for idx, k := range ks {
			elem, err := txn.Get(bucket_, k)
			if err == nutsdb.ErrKeyNotFound{
				if err = fn(idx, k, nil); err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
			if err := fn(idx, k, elem.Value); err != nil {
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
	fn := func(txn *nutsdb.Tx) error {
		if bucket != nil {
			if len(*bucket) == 0 {
				return fmt.Errorf("empty bucket name")
			}
		} else {
			bucket = &empty_bucket_name
		}
		bucket_ := *bucket
		for i, k := range ks {
			elem, err := txn.Get(bucket_, k)
			if err == nutsdb.ErrKeyNotFound {
				out[i] = nil
				continue
			}
			if err != nil {
				out = nil
				return err
			}

			out[i] = elem.Value
			if del {
				if err := txn.Delete(bucket_, k); err != nil{
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
	return db.db.Update(func(txn *nutsdb.Tx) error {
		for _, k := range ks{
			if err := txn.Delete(empty_bucket_name, k); err != nil{
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
	return db.setIFs(&bucket, items, fn)
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
	if bucket != nil{
		if len(*bucket) == 0 {
			return fmt.Errorf("empty bucket name")
		}
		if db.buckets[*bucket] == nil {
			return nil
		}
	} else {
		bucket = &empty_bucket_name
	}
	return db.db.View(func(txn *nutsdb.Tx) error {
		idx := -1

		es, err := txn.GetAll(*bucket)
		if err != nil {
			return err
		}

		for _, e := range es {
			err := fn(idx, e.Key, e.Value)
			if err != nil {
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
		
		err := db.db.Update(func(txn *nutsdb.Tx) error {
			return txn.DeleteBucket(nutsdb.DataStructureNone, *bucket)
		})
		if err != nil {
			delete(db.buckets, *bucket)
		}
	} else {
		bucket = &empty_bucket_name
		return db.db.Update(func(txn *nutsdb.Tx) error {
			return txn.DeleteBucket(nutsdb.DataStructureNone, *bucket)
	  })
  }

	return nil
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
	return db.db.Update(func(txn *nutsdb.Tx) error {
		return txn.Delete(bucket, k)
	})
}

func (db *DB)BDels(bucket string, ks [][]byte)(error) {
	if len(bucket) == 0 {
    return fmt.Errorf("empty bucket name")
	}

	return db.db.Update(func(txn *nutsdb.Tx) error {
		for _, k := range ks {
			if err := txn.Delete(bucket, k); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB)rmLocalFiles() error {

	dir, err := ioutil.ReadDir(db.opts.Dir)
	if err != nil {
		return err
	}
  for _, d := range dir {
    err = os.RemoveAll(path.Join([]string{db.opts.Dir, d.Name()}...))
		if err != nil {
			return err
		}
  }

	return nil
}

func(db *DB)Truncate() (error) {
	err := db.db.Close()
	if err != nil {
		return err
	}

	err = db.rmLocalFiles()
	if err != nil {
		return err
	}

	db.db, err = nutsdb.Open(*db.opts)
	if err != nil {
		return err
	}

	return nil
}