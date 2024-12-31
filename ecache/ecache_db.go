package ecache

import (
	"fmt"
	"time"

	"github.com/ziyht/eden_go/ecache/driver"
	"github.com/ziyht/eden_go/eutils/ptr"
)

type db struct {
  dsn string
	db  driver.DB
}

func newDB(opts *DBOpts) (*db, error) {
	if opts.Dsn == "" {
		opts.Dsn = GenDsn(opts.Driver, opts.Dir, opts.Params)
	}

	db_, err := driver.OpenDsn(opts.Dsn)
	if err != nil {
		return nil, err
	}

	if db_ == nil {
		return nil, fmt.Errorf("invalid returned db(nil) checked from current driver in dsn(%s)", opts.Dsn)
	}

	return &db{dsn: opts.Dsn, db: db_}, nil
}

func (db *db)setVal(pre []byte, key []byte, val Val, ttl ...time.Duration)error{
	if val.Type() == Nil || val.Type() >= VT_MAX{
		return fmt.Errorf("invalid value type(%s) to set to DB", val.__typeStr())
	}

	return db.db.Update(func(tx driver.TX)error{
		return tx.Set(pre, key, val.marshal(), ttl...)
	})
}

func (db *db)setAny(pre []byte, key any, val any, ttl ...time.Duration) error {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}

	raw, err := NewVal(val)
	if err != nil {
		return err
	}

	defer recycleVal(raw)
	return db.db.Update(func(tx driver.TX)error{
		return tx.Set(pre, k, raw.marshal(), ttl...)
	})
}

type __db_val_getter func(vals [][]byte, i int)([]byte)
type __db_ttl_getter func(ttls []time.Duration, i int)([]time.Duration)

func (db *db)sets(prefix []byte, keys [][]byte, vals [][]byte, ttls ... time.Duration) error {
	lk := len(keys); 
	if lk == 0 {
		return nil
	}

	lv := len(vals); 
	ld := len(ttls);
	var _val_getter __db_val_getter
	var _ttl_getter __db_ttl_getter

	switch lv {
		case 1 : _val_getter = func (vals [][]byte, i int)([]byte){ return vals[0] }
		case lk: _val_getter = func (vals [][]byte, i int)([]byte){ return vals[i] }
		default: return fmt.Errorf("the len for keys(%d) and values(%d) are not match, and len of values is not 1", lk, lv)
	}

	switch ld {
		case 0 : _ttl_getter = func(ttls []time.Duration, i int)([]time.Duration){ return nil  };
		case 1 : _ttl_getter = func(ttls []time.Duration, i int)([]time.Duration){ return ttls };
		case lk: ttl_cache := make([]time.Duration, 1)
						_ttl_getter = func(ttls []time.Duration, i int)([]time.Duration){ ttl_cache[0] = ttls[i]; return ttl_cache };
		default: return fmt.Errorf("the len for key(%d) and ttls(%d) are not match, and len of ttls is not 1 or 0", lk, ld)
	}

	for i := 0; i < lk;  {
		if err := db.db.Update(func(tx driver.TX)error{
			cnt := 0
			var val Val;
			for ; i < lk && cnt < 1000; i++ {

				bin := _val_getter(vals, i)
				val.setBytes(bin)
				if err := tx.Set(prefix, keys[i], val.marshal(), _ttl_getter(ttls, i)...); err != nil {
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

func (db *db)setsAny(prefix []byte, keys any, vals any, ttls ... time.Duration) error {
	ks, err := toBytesArr(keys)
	if err != nil {
		return fmt.Errorf("invalid keys input: %s", err)
	}
	vs, err := toBytesArr(vals)
	if err != nil {
		return fmt.Errorf("invalid keys input: %s", err)
	}
	return db.sets(prefix, ks, vs, ttls...)
}

func (db *db)getAny(prefix []byte, key any, del ...bool)(val Val, err error) {
	k, err := toBytesKey(key)
	if err != nil {
		return
	}

	err = db.db.View(func(tx driver.TX)error{
		bin, _, err := tx.Get(prefix, k, del...)
		val.unmarshal(bin)
		return err
	})

	return
}

func (db *db)getAnyEx(prefix []byte, key any, del ...bool)(val Val, expiresAt uint64, err error) {
	k, err := toBytesKey(key)
	if err != nil {
		return
	}

	err = db.db.View(func(tx driver.TX)error{
		val.d, expiresAt, err = tx.Get(prefix, k, del...)
		val.unmarshal(val.d)
		return err
	})

	return
}

func (db *db)getBytesExt(prefix []byte, key []byte, del ...bool)(val []byte, expiresAt uint64, err error) {
	db.db.View(func(tx driver.TX)error{
		val, expiresAt, err = tx.Get(prefix, key, del...)
		return err
	})

	return
}

func (db *db)gets(prefix []byte, keys [][]byte, del ...bool)(vals []Val, err error) {
	err = db.db.View(func(tx driver.TX)error{
		var val Val
		for _, key := range keys {
			val.d, _, err = tx.Get(prefix, key, del...)
			if err != nil {
				return err
			}
			val.unmarshal(val.d)
			vals = append(vals, val)
		}

		return nil
	})

	return
}

func (db *db)getsAny(prefix []byte, keys any, del ...bool)(vals []Val, err error) {
	err = db.db.View(func(tx driver.TX)error{
		var val Val
		switch k := keys.(type) {
			case string  : bin, _, err := tx.Get(prefix, ptr.StringToBytes(k), del...); if err != nil { return err }; val.unmarshal(bin); vals = append(vals, val)
			case []byte  : bin, _, err := tx.Get(prefix,                   k , del...); if err != nil { return err }; val.unmarshal(bin); vals = append(vals, val)
			case []string: for _, tk := range k { bin, _, err := tx.Get(prefix, ptr.StringToBytes(tk), del...); if err != nil { return err }; val.unmarshal(bin); vals = append(vals, val)  }
			case [][]byte: for _, tk := range k { bin, _, err := tx.Get(prefix,                   tk , del...); if err != nil { return err }; val.unmarshal(bin); vals = append(vals, val)  }
			default      : return fmt.Errorf("invalid key type, only support: string, []string, []byte and [][]byte")
		}

		return nil
	})

	return
}

func (db *db)getAll(prefix []byte, restore... bool)(keys [][]byte, vals []Val, err error) {
	err = db.db.View(func(tx driver.TX)error{
		var val Val
		return tx.Iterate(prefix, func(_ int, key []byte, val_ []byte, _ uint64)error{
			val.unmarshal(val_)

			keys = append(keys, key)
			vals = append(vals, val)
			return nil
		})
	})

	return
}

func (db *db)del(prefix []byte, key []byte)(err error) {
	return db.db.Update(func(tx driver.TX)error{
		return tx.Del(prefix, key)
	})
}

func (db *db)delAny(prefix []byte, key any)(err error) {
	k, err := toBytesKey(key)
	if err != nil {
		return err
	}
	return db.db.Update(func(tx driver.TX)error{
		return tx.Del(prefix, k)
	})
}

func (db *db)dels(prefix []byte, keys ...[]byte)(err error) {
	return db.db.Update(func(tx driver.TX)error{
		for _, key := range keys {
			if err = tx.Del(prefix, key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *db)delsAny(prefix []byte, keys ...any) error {
	return db.db.Update(func(tx driver.TX)error{
		for idx, anyk := range keys {
			switch k := anyk.(type) {
				case string  : if err := tx.Del(prefix, []byte(k)); err != nil { return err }
				case []byte  : if err := tx.Del(prefix,        k ); err != nil { return err }
				case []string: for _, tk := range k { if err := tx.Del(prefix, []byte(tk)); err != nil { return err }  }
				case [][]byte: for _, tk := range k { if err := tx.Del(prefix,        tk ); err != nil { return err }  }
				default: return fmt.Errorf("invalid key type(%t) at idx(:%d), only support: string, []string, []byte and [][]byte", k, idx)
			}
		}
		return nil
	})
}

func (db *db)setObjs(prefix []byte, objs []any, fn func(int, any)(key []byte, val any, ttl time.Duration))(err error) {
	len := len(objs)
	for i := 0; i < len; {
		if err = db.db.Update(func(tx driver.TX)error{
			cnt := 0
			for ; i < len && cnt < 1000; i++ {
				key, val_, ttl := fn(i, objs[i])

				val, err := NewVal(val_)
				if err != nil {
					return err
				}

				err = tx.Set(prefix, key, val.marshal(), ttl)
				if err != nil {
					return err
				}
				recycleVal(val)
				cnt += 1
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *db)doForAll(prefix []byte, fn func(idx int, key []byte, val Val) error) (err error) {
	return db.db.View(func(tx driver.TX)error{
		var val Val
		return tx.Iterate(prefix, func(idx int, key []byte, bin []byte, _ uint64)error{
			val.unmarshal(bin)
			return fn(idx, key, val)
		})
	})
}

func (db *db)doForAllEx(prefix []byte, fn func(idx int, key []byte, val Val, expiresAt uint64) error) (err error) {
	return db.db.View(func(tx driver.TX)error{
		var val Val
		return tx.Iterate(prefix, func(idx int, key []byte, bin []byte, expiresAt uint64)error{
			val.unmarshal(bin)
			return fn(idx,  key, val, expiresAt)
		})
	})
}

func (db *db)doForKeys(prefix []byte, keys [][]byte, fn func(idx int, key []byte, val Val) error) (err error) {
	return db.db.View(func(tx driver.TX)error{
		var val Val

		for i, key := range keys {
			bin, _, err := tx.Get(prefix, key)
			if err != nil {
				return err 
			}

			val.unmarshal(bin)
			if err = fn(i, key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *db)doForKeysAny(prefix []byte, keys any, fn func(idx int, key []byte, val Val) error) (err error) {
	ks, err := toBytesArr(keys)
	if err != nil {
		return fmt.Errorf("invalid type(%t) of keys: %s", keys, err)
	}

	var val Val
	return db.db.View(func(tx driver.TX)error{
		for i, key := range ks {
			bin, _, err := tx.Get(prefix, key)
			if err != nil {
				return err 
			}

			val.unmarshal(bin)
			if err = fn(i, key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *db)close() error {
	return db.db.Close()
}

func (db *db)truncate() error {
	return db.db.Truncate()
}




