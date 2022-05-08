package ecache

import (
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type ECacheL1 struct {
	bucket string
	c      *ECache
	cache  cache.Cache
	max    time.Duration
	rwmu   sync.RWMutex
	safe   bool
}

func (l1 *ECacheL1)setToMem(k string, i Item, ttl ...time.Duration) {
	ttl_to_set := time.Duration(0)
	if len(ttl) == 0 {
		ttl_to_set = i.TTL()
	} else {
		ttl_to_set = ttl[0]		
	}

	if ttl_to_set > l1.max || ttl_to_set == 0{
		ttl_to_set = l1.max
	}

	l1.cache.Set(k, i, ttl_to_set)
}

func (l1 *ECacheL1)rmFromMem(k string) {
	l1.cache.Delete(k)
}

func (l1 *ECacheL1)rLock()   { if l1.safe {l1.rwmu.RLock()}}
func (l1 *ECacheL1)rUnlock() { if l1.safe {l1.rwmu.RUnlock()}}
func (l1 *ECacheL1)wLock()   { if l1.safe {l1.rwmu.Lock()}}
func (l1 *ECacheL1)wUnlock() { if l1.safe {l1.rwmu.Unlock()}}

func (l1 *ECacheL1)setItemToDBCache(k string, i Item, ttl ...time.Duration) error {
	if l1.bucket == "" {
		return l1.c.SetItem([]byte(k), i, ttl...)
	} 
	return l1.c.BSetItem(l1.bucket, []byte(k), i, ttl...)
}
func (l1 *ECacheL1)getItemFromDBCache(k string, i Item) (Item, error) {
	if l1.bucket == "" {
		return l1.c.GetItem([]byte(k), i)
	}
	return l1.c.BGetItem(l1.bucket, []byte(k), i)
}

// GetL1 - Get a l1 cache, this cache has a internal mem cache
// bucket can be empty
// the bucket in same name has the same l1 cache instance, include empty bucket
// if already exist, return the earlier instance directly
func (c *ECache)GetL1(bucket string, max time.Duration, safe ...bool) *ECacheL1{
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.l1s == nil {
		c.l1s = make(map[string]*ECacheL1)
	}

	l1 := c.l1s[bucket]
	if l1 == nil{
		l1 = &ECacheL1{
			bucket: bucket,
			c     : c,
			cache : *cache.New(max, time.Hour),
			max   : max,
		}

		if len(safe) > 1 {
			l1.safe = safe[0]
		}

		c.l1s[bucket] = l1
	}

	return l1
}

func (l1 *ECacheL1)L1Count() int {
	return l1.cache.ItemCount()
}

// SetItem
// the ttl in params has higher priority than item.TTL()
func (l1 *ECacheL1)SetItem(k string, i Item, ttl ...time.Duration) (err error) {
	l1.wLock()
	defer l1.wUnlock()

	if err = l1.setItemToDBCache(k, i, ttl...); err != nil {
		return err
	}

	l1.setToMem(k, i, ttl...)
	return nil
}

func (l1 *ECacheL1)GetItem(k string, i Item) (out Item, err error) {
	l1.rLock()
	defer l1.rUnlock()

	i_, exist := l1.cache.Get(k)
	if exist {

		out, ok := i_.(Item)
		if !ok {
			return nil, ErrConvertFailed
		}

		return out, nil
	}

	out, err = l1.getItemFromDBCache(k, i)
	if err == nil && out != nil{
		l1.setToMem(k, out, out.TTL())
	}
	return
}





