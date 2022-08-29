package ecache

import (
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/ziyht/eden_go/ecache/driver"
)

type Item interface {
	Marshal()([]byte, error)  // do marshal things
	Unmarshal([]byte) error   // ummarshal data to self
}

type ItemRegion [T Item] struct {
  db       *db
	meta     rMeta
	ttl      time.Duration
	mem      *MemCache[T]
	Metrics  *ristretto.Metrics
}

func newItemRegion[T Item](db *db, ks []string) (*ItemRegion[T]) {
	r := &ItemRegion[T]{db: db}
	r.meta.initItem(ks)	

	return r
}

func (r *ItemRegion[T])EnableMemCache(maxCount int64, maxTTL time.Duration) {
	if r.mem != nil {
		return 
	}

	opts := MemCacheOpts[T]{
		CountersNum: maxCount * 10,
		MaxCost    : maxCount,
		TTL        : r.ttl,
		MaxTTL     : maxTTL,
		AutoReRent : true,
		IgnoreInternalCost: true,
		Statistics : true,
	}
	r.mem = newMemCache(opts)
	r.Metrics = r.mem.Metrics
}

func (r *ItemRegion[T])SetDefaultTTL(ttl time.Duration){
	r.ttl = ttl
}

func (r *ItemRegion[T])setToMem(k []byte, i T, cost int64, ttl ...time.Duration) {
	if r.mem == nil {
		return
	}

	r.mem.SetEx(k, i, 1, ttl...)
}

func (r *ItemRegion[T])getFromMem(k []byte, del...bool) (out T, ok bool) {
	if r.mem == nil {
		return 
	}

	val, ok := r.mem.Get(k)
	if ok{
		if len(del) > 0 && del[0] {
			r.mem.Del(k)
		}

		return val, true
	}

	return 
}

func (r *ItemRegion[T])Del(key []byte) error {
	r.mem.Del(key)
	return r.db.del(r.meta.kpre, key)
}

func (r *ItemRegion[T])ADel(key any) error {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}

	r.mem.Del(key)
	return r.db.del(r.meta.kpre, k)
}

func (r *ItemRegion[T])Set(key []byte, item T, ttl ...time.Duration) error {
	val, err := item.Marshal(); 
	if err != nil {
		return err
	}
	
	valid_ttl := r.ttl
	if len(ttl) > 0 && ttl[0] >= 0 {
		valid_ttl = ttl[0]
	}

	r.setToMem(key, item, 1, valid_ttl)
	return r.db.set(r.meta.kpre, key, val, valid_ttl)
}

func (r *ItemRegion[T])ASet(key any, item T, ttl ...time.Duration) error {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}
	return r.Set(k, item, ttl...)
}

func (r *ItemRegion[T])Get(key []byte, new func() T, del...bool)(out T, err error) {
	c, ok := r.getFromMem(key, del...)
	if ok {
		if len(del) > 0 && del[0] {
			r.Del(key)
		}

		return c, nil
	}

	i := new()

	data, expire, err := r.db.getExt(r.meta.kpre, key, del...)
	if err != nil {
		return
	}

	if data == nil {
		return
	}

	if err = i.Unmarshal(data); err != nil {
		return 
	}

	if len(del) > 0 && del[0] {
		r.db.del(r.meta.kpre, key)
	} else if expire == 0 {
		r.setToMem(key, i, 1, time.Duration(0))
	} else {
		r.setToMem(key, i, 1, time.Until(time.Unix(int64(expire), 0)))
	}

	return i, nil
}

func (r *ItemRegion[T])AGet(key any, new func() T, del...bool)(out T, err error){
	k, err := toBytesKey(key)
	if err != nil {
		return 
	}
	return r.Get(k, new, del...)
}

// Gets
// - skipErrs_Del_RetainNil: this is a three-bool-value to set options
//   skipErrs : if is true, it will continue when err occurs in Unmarshal operations
//   Del      : if is true, the keys will be deleted after all the operations
//   RetainNil: if is true, the nil value which created by err and not_found will be retain in the results
func (r *ItemRegion[T])Gets(keys [][]byte, new func() T, skipErrs_Del_RetainNil...bool)(items []T, err error){
	if len(keys) == 0 {
		return nil, err
	}

	skipErr   := len(skipErrs_Del_RetainNil) > 0 && skipErrs_Del_RetainNil[0]
	retainNil := len(skipErrs_Del_RetainNil) > 2 && skipErrs_Del_RetainNil[2]

	if err := r.db.db.View(func(tx driver.TX)error{
		for _, key := range keys {
			i, ok := r.getFromMem(key)
			if ok {
				items = append(items, i)
				continue
			}

			val, _, err := tx.Get(r.meta.kpre, key)
			if err == nil {
				i2 := new()
				err = i2.Unmarshal(val)
				if err == nil {
					items = append(items, i2)
					continue
				}
			}
			
			if skipErr {
				if retainNil{
					items = append(items, i)
				}

				continue
			}

			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if len(skipErrs_Del_RetainNil) > 1 && skipErrs_Del_RetainNil[1] {
		r.mem.Dels(keys)
		r.db.dels(r.meta.kpre, keys...)
		r.mem.Wait()
	} 

	return 
}

func (r *ItemRegion[T])AGets(keys any, new func() T, skipErrs_Del_RetainNil...bool)(items []T, err error){
	ks, err := toBytesArr(keys)
	if err != nil {
		return nil, fmt.Errorf("invalid type(%t) of input keys: %s", keys, err)
	}
	return r.Gets(ks, new, skipErrs_Del_RetainNil...)
}

func (r *ItemRegion[T])GetAll(new func() T, skipErrs... bool) ([]T, error) {
	var items []T
	skipErr := len(skipErrs) > 0 && skipErrs[0]
	err := r.db.doForAllEx(r.meta.kpre, func(idx int, key []byte, val []byte, expiresAt uint64) error{

		i := new()
		err := i.Unmarshal(val)
		if err == nil {
			items = append(items, i)
			r.setToMem(key, i, 1, time.Until(time.Unix(int64(expiresAt), 0)))
		} else if !skipErr{
			return fmt.Errorf("do Unmarshal failed for key '%s': %s", key, err)
		}

		return nil
	})

	return items, err
}

func (r *ItemRegion[T])ReloadItems(new func() T)(int, error) {
	if r.mem == nil {
		return 0, fmt.Errorf("internal mem")
	}

	cnt := 0
	err := r.db.doForAllEx(r.meta.kpre, func(idx int, key []byte, val []byte, expiresAt uint64) error{

		i := new()
		err := i.Unmarshal(val)
		if err == nil {
		cnt += 1
			r.setToMem(key, i, 1, time.Until(time.Unix(int64(expiresAt), 0)))
		}

		return nil
	})

	return cnt, err
}