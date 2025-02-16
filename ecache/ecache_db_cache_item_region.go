package ecache

import (
	"fmt"
	"time"

	//"github.com/dgraph-io/ristretto"
	"github.com/ziyht/eden_go/ecache/driver"
)

type Item interface {
	Marshal()([]byte, error)  // do marshal things
	Unmarshal([]byte) error   // ummarshal data to self
}

type ItemRegion[V Item] struct {
  db       *db
	meta     rMeta
	ttl      time.Duration
	mem      *MemCache[[]byte, V]
	Metrics  *Metrics
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
		DfTTL      : r.ttl,
		MaxTTL     : maxTTL,
		AutoReRent : true,
		IgnoreInternalCost: true,
		Statistics : true,
	}
	r.mem = newMemCache[[]byte, T](opts)
	r.Metrics = r.mem.Metrics
}

func (r *ItemRegion[T])SetDefaultTTL(ttl time.Duration){
	r.ttl = ttl
}

func (r *ItemRegion[T])setToMem(k []byte, v T, cost int64, ttl ...time.Duration) {
	if r.mem == nil {
		return
	}

	r.mem.SetEx(k, v, cost, ttl...)
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

// key and val can only be string or []byte
func (r *ItemRegion[T])Del(key []byte) error {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}

	r.mem.Del(key)
	return r.db.del(r.meta.kpre, k)
}

// key and val can only be string or []byte
func (r *ItemRegion[T])Set(key any, item T, ttl ...time.Duration) error {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}

	var raw Val
	err = raw.setItem(item)
	if err != nil {
		return err
	}
	
	valid_ttl := r.ttl
	if len(ttl) > 0 && ttl[0] >= 0 {
		valid_ttl = ttl[0]
	}

	r.setToMem(k, item, 1, valid_ttl)
	return r.db.setVal(r.meta.kpre, k, raw, valid_ttl)
}

func (r *ItemRegion[T])__valToItem(val Val, new func() T)(out T, err error){
  i := new()

	err = val.Error()
	if err != nil {
		return
	}

	if val.Type() != ITEM{
		err = fmt.Errorf("invalid Raw type(%s)", val.__typeStr())
		return
	}

	if err = i.Unmarshal(val.d); err != nil {
		return 
	}

	return i, nil
}

// key and val can only be string or []byte
func (r *ItemRegion[T])Get(key any, new func() T, del...bool)(out T, err error) {
	k, err := toBytesKey(key)
	if err != nil {
		return 
	}

	c, ok := r.getFromMem(k, del...)
	if ok {
		if len(del) > 0 && del[0] {
			r.Del(k)
		}

		return c, nil
	}

	bin, expire, err := r.db.getBytesExt(r.meta.kpre, k, del...)
	if err != nil {
		return
	}

	if bin == nil {
		return
	}

	var val Val; val.unmarshal(bin)
	i, err := r.__valToItem(val, new)
	if err != nil {
		return
	}

	if len(del) > 0 && del[0] {
		r.db.del(r.meta.kpre, k)
	} else if expire == 0 {
		r.setToMem(k, i, 1, time.Duration(0))
	} else {
		r.setToMem(k, i, 1, time.Until(time.Unix(int64(expire), 0)))
	}

	return i, nil
}

// Gets
//   keys: can only be string, []string, []byte or [][]byte
//   skipErrs_Del_RetainNil: this is a three-bool-value options:
//     skipErrs : if is true, it will continue when err occurs in Unmarshal operations
//     Del      : if is true, the keys will be deleted after all the operations
//     RetainNil: if is true, the nil value which created by err and not_found will be retain in the results
func (r *ItemRegion[T])Gets(keys [][]byte, new func() T, skipErrs_Del_RetainNil...bool)(items []T, err error){
	// ks, err := toBytesArr(keys)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid type(%t) of input keys: %s", keys, err)
	// }

	ks := keys

	if len(ks) == 0 {
		return nil, err
	}

	skipErr   := len(skipErrs_Del_RetainNil) > 0 && skipErrs_Del_RetainNil[0]
	retainNil := len(skipErrs_Del_RetainNil) > 2 && skipErrs_Del_RetainNil[2]

	if err := r.db.db.View(func(tx driver.TX)error{
		for _, key := range ks {
			i, ok := r.getFromMem(key)
			if ok {
				items = append(items, i)
				continue
			}

			bin, _, err := tx.Get(r.meta.kpre, key)
			if err == nil {
				var val Val; val.unmarshal(bin)
				i2, err := r.__valToItem(val, new)
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
		for _, key := range ks {
			r.mem.Del(key)
		}
		r.db.dels(r.meta.kpre, ks...)
		r.mem.Wait()
	} 

	return 
}

func (r *ItemRegion[T])GetAll(new func() T, skipErrs... bool) ([]T, error) {
	var items []T
	skipErr := len(skipErrs) > 0 && skipErrs[0]
	err := r.db.doForAllEx(r.meta.kpre, func(idx int, key []byte, val Val, expiresAt uint64) error{

		i, err := r.__valToItem(val, new)
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
		return 0, fmt.Errorf("internal memcache not enabled")
	}

	cnt := 0
	err := r.db.doForAllEx(r.meta.kpre, func(idx int, key []byte, val Val, expiresAt uint64) error{

		i, err := r.__valToItem(val, new)
		if err == nil {
		cnt += 1
			r.setToMem(key, i, 1, time.Until(time.Unix(int64(expiresAt), 0)))
		}

		return nil
	})

	return cnt, err
}