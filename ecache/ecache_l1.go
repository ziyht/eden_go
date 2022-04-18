package ecache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

type ECacheL1 struct {
	bucket string
	c      *ECache
	l1     cache.Cache
	max    time.Duration
}

func (c *ECache)GetL1(bucket string, max time.Duration) *ECacheL1{
	c.mu.Lock()
	defer c.mu.Unlock()

	l1 := c.l1s[bucket]
	if l1 == nil{
		l1 = &ECacheL1{
			bucket: bucket,
			c     : c,
			l1    : *cache.New(max, time.Hour),
			max   : max,
		}
		c.l1s[bucket] = l1
	}

	return l1
}

func (l1 *ECacheL1)setToL1(k string, i Item, ttl ...time.Duration) {
	if len(ttl) == 0 {
		l1.l1.SetDefault(k, i)
	} else {
		ttl_set := ttl[0]
		if ttl_set > l1.max {
			ttl_set = l1.max
		}
		l1.l1.Set(k, i, ttl_set)
	}
}

func (l1 *ECacheL1)SetI(k string, i Item, ttl ...time.Duration) (err error) {

	if err = l1.c.BSetI(l1.bucket, []byte(k), i, ttl...); err != nil {
		return err
	}
	l1.setToL1(k, i, ttl...)
	return nil
}

func (l1 *ECacheL1)GetI(k string, i Item) (any, error) {
	i_, exist := l1.l1.Get(k)
	if exist {
		return i_, nil
	}

	return l1.c.BGetI("", []byte(k), i)
}