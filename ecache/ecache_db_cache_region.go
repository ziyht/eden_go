package ecache

import (
	"time"
)

type Region struct {
  db       *db
	meta     rMeta
	ttl      time.Duration   // not used now
}

func newRegion(db *db, ks []string) (*Region) {
	out := &Region{db: db}
	out.meta.init(ks)
	return out
}

func (r *Region)SubRegion(ks ...string) (*Region){
	if len(ks) == 0 {
		r = &Region{db: r.db, meta: r.meta}
		return r
	}

	out := &Region{db: r.db}
	out.meta = r.meta.genSubMeta(ks)
	return out
}

func (r *Region)ToItemRegion()(*ItemRegion[Item]) {
	return newItemRegion[Item](r.db, r.meta.keys)
}

func (r *Region)SetDefaultTTL(ttl time.Duration) {
	r.ttl = ttl
}

func (r *Region)Set(key []byte, val []byte, ttl ...time.Duration) error {
	if len(ttl) > 0 {
		return r.db.set(r.meta.kpre, key, val, ttl...)
	}

	return r.db.set(r.meta.kpre, key, val, r.ttl)
}

// key and val can only be string or []byte
func (r *Region)ASet(key any, val any, ttl ...time.Duration) error {
	if len(ttl) > 0 {
		return r.db.setAny(r.meta.kpre, key, val, ttl...)
	}
	return r.db.setAny(r.meta.kpre, key, val, r.ttl)
}

func (r *Region)Sets(keys [][]byte, vals [][]byte, ttls ...time.Duration) error {
	if len(ttls) > 0 {
		return r.db.sets(r.meta.kpre, keys, vals, ttls...)
	}

	return r.db.sets(r.meta.kpre, keys, vals, r.ttl)
}

// key and val can only be string, []string, []byte or [][]byte
func (r *Region)ASets(keys any, vals any, ttls ...time.Duration) error {
	if len(ttls) > 0 {
		return r.db.setsAny(r.meta.kpre, keys, vals, ttls...)
	}

	return r.db.setsAny(r.meta.kpre, keys, vals, r.ttl)
}

func (r *Region)SetObjs(items []any, fn func(int, any)(k[]byte, v[]byte, ttl time.Duration))error{
	return r.db.setObjs(r.meta.kpre, items, fn)
}

func (r *Region)Get(key []byte, del ...bool)([]byte, error){
	return r.db.get(r.meta.kpre, key, del...)
}

// key can only be string or []byte
func (r *Region)AGet(key any, del ...bool)([]byte, error){
	return r.db.getAny(r.meta.kpre, key, del...)
}

func (r *Region)Gets(keys [][]byte, del ...bool)([][]byte, error){
	return r.db.gets(r.meta.kpre, keys, del...)
}

// key can only be string, []strng, []byte or [][]byte
func (r *Region)AGets(keys any, del ...bool)([][]byte, error){
	return r.db.getsAny(r.meta.kpre, keys, del...)
}

// GetAll - returns all keys and values in this region
func (r *Region)GetAll()([][]byte, [][]byte, error){
	return r.db.getAll(r.meta.kpre)
}

func (r *Region)Del(key []byte)(error){
	return r.db.del(r.meta.kpre, key)
}

// key can only be string or []byte
func (r *Region)ADel(key any)(error){
	return r.db.delAny(r.meta.kpre, key)
}

func (r *Region)Dels(keys... []byte)(error){
	return r.db.dels(r.meta.kpre, keys...)
}

// key can only be string, []strng, []byte or [][]byte
func (r *Region)ADels(keys ...any) (error) {
	return r.db.delsAny(r.meta.kpre, keys...)
}

func (r *Region)DoForAll(fn func(idx int, key []byte, val []byte) error)error{
	return r.db.doForAll(r.meta.kpre, fn)
}

func (r *Region)DoForKeys(keys [][]byte, fn func(idx int, key []byte, val []byte) error)error{
	return r.db.doForKeys(r.meta.kpre, keys, fn)
}

func (r *Region)ADoForKeys(keys any, fn func(idx int, key []byte, val []byte) error)error{
	return r.db.doForKeysAny(r.meta.kpre, keys, fn)
}

func (r *Region)Truncate(/*including_subs ...bool*/) error {
	// if len(including_subs) > 0 && including_subs[0] {
	// 	// TODO
	// 	return fmt.Errorf("todo")
	// }

	return r.db.db.DropPrefix(r.meta.kpre)
}
