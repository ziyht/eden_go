package ecache

import (
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
)

type MemCache[T any] struct {
  c                *ristretto.Cache
	rerentSkipThresh time.Duration      // default MaxTTL * 0.9
	opts             MemCacheOpts[T]
	Metrics          *ristretto.Metrics
	add              atomic.Int32
}

type memCacheItem[T any] struct {
  cost          int
  needExpiredAt int64           // a millisecond timestamp, this is used for auto rerent, 0 means auto rerent always
	Value         T
}

const memCacheItemCost = int(unsafe.Sizeof(memCacheItem[any]{}))

type MemCacheOpts[T any] struct {
	MaxCost            int64             // default:   10M, each item has a cost of memory, this value defines the max cost of current cache instance
	TTL                time.Duration     // default:     0, the default TTL of items when not passed ttl in params in Set functions, 0 means the item will never expired by TTL policy.
  MaxTTL             time.Duration     // default:     0, the max TTL of all items, the ttl will set to it if the passed ttl is greater than MaxTTL
	AutoReRent         bool              // default: false, the ttl of the item will be automatically increasing to MaxTTL(if is set) after access it
	OnEvict            func(T)           // default:   nil, this will be called on each item evict 
	OnCost             func(T) int       // default:   nil, this will be called to get real cost of item when the input cost is 0
	IgnoreInternalCost bool              // default: false, IgnoreInternalCost set to true indicates to the cache that the cost of internally storing the value should be ignored.
	CountersNum        int64             // default:  4096, CountersNum determines the number of counters (keys) to keep that hold access frequency information used by internal policy, not the count limit for items. It's generally a good idea to have more counters(10x) than the max cache capacity, as this will improve eviction accuracy and subsequent hit ratios.
	BufferItems        int64             // default:    64, ristretto: BufferItems is the size of the Get buffers. The best value we've found for this is 64.
	Statistics         bool              // default: false, do Statistics interval or not, 
}

func __cost_of_any(i any) int {
	if i == nil {
		return 0
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Struct : return int(unsafe.Sizeof(v))
	case reflect.Ptr, reflect.Interface:
			p := (*[]byte)(unsafe.Pointer(v.Pointer()))
			if p == nil {
					return 0
			}
			return int(unsafe.Sizeof(v.Elem()))
	case reflect.String, reflect.Array, reflect.Chan, reflect.Slice:
			return v.Len()
	case reflect.Bool, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Int:
			return int(v.Type().Size())
	default :
		return int(unsafe.Sizeof(v))
	}
}

// 
func (c *MemCache[T])__cost_eval_func(i any) (int64) {
	item, ok := i.(*memCacheItem[T])
	if !ok {
		if !c.opts.IgnoreInternalCost{
			return int64(memCacheItemCost)
		}
		return 0
	}

	cost := item.cost
	if cost == 0 {
		cost = c.opts.OnCost(item.Value)
		if cost < 0 {
			cost = 0
		}
	}
  
	if !c.opts.IgnoreInternalCost {
		cost += memCacheItemCost
	}

	item.cost = cost

	return int64(cost)
}

func newMemCache[T any](opts MemCacheOpts[T]) *MemCache[T] {

	if opts.CountersNum <= 0 {
		opts.CountersNum = 4096
	}
	if opts.MaxCost <= 0 {
		opts.MaxCost = 10 * 1024 * 1024          // 1M
	}
	if opts.BufferItems <= 0 {
		opts.BufferItems = 64
	}
	if opts.OnCost == nil {
		opts.OnCost = func(i T)int{
			return __cost_of_any(i)
		}
	}

	out := &MemCache[T]{
		opts     : opts,
	}
	out.__cal_rerent_ttl()

	config := &ristretto.Config{
		NumCounters       : opts.CountersNum,
		MaxCost           : opts.MaxCost,
		BufferItems       : opts.BufferItems,
		IgnoreInternalCost: opts.IgnoreInternalCost,
		Cost              : out.__cost_eval_func,
		Metrics           : opts.Statistics,
	}

	if opts.OnEvict != nil {
		config.OnEvict = func(i *ristretto.Item){
			item, _ := i.Value.(*memCacheItem[T])
			opts.OnEvict(item.Value)
		}
	}

	out.c, _ = ristretto.NewCache(config)
	out.Metrics = out.c.Metrics
	if out.c == nil {
		return nil
	}

	return out
}

// the internal cache is not synchronized with Set funcs
// you should call Wait() after call Set funcs
func (c *MemCache[T])Wait() {
	c.c.Wait()
}

func (c *MemCache[T])__cal_rerent_ttl(){
	c.rerentSkipThresh = time.Duration(float64(c.opts.MaxTTL) * 0.9)
}

func (c *MemCache[T])__make_item(value T, cost int, ttl ...time.Duration)(*memCacheItem[T], time.Duration){
	ttl_set := c.opts.TTL
	if len(ttl) > 0 {
		ttl_set = ttl[0]
	}
	
	if cost != 0 && !c.opts.IgnoreInternalCost{
		cost += memCacheItemCost
	}

	//
	// if ttl_set is 0, means this item should never been expired
	// so max_ttl take control:
	//   1. if max_ttl is  0, return max_ttl means return 0, the item will never expirrd
	//   2. if max_ttl is >0, return max_ttl means means the item will expire after max_ttl
	// 
	// here we not set needExpiredAt(==0), that means this item will auto rerent always on each access if max_ttl > 0
	// 
	if ttl_set == 0 {
		return &memCacheItem[T]{cost: cost, Value: value}, c.opts.MaxTTL
	}

	// ttl_set > 0, means this item should expire after ttl_set for whatever max_ttl is
	item := &memCacheItem[T]{cost: cost, Value: value, needExpiredAt: time.Now().Add(ttl_set).UnixMilli()}
	// if max_ttl is valid, so the current ttl should not bigger than it
	if c.opts.MaxTTL> 0 && c.opts.MaxTTL < ttl_set {
		ttl_set = c.opts.MaxTTL
	}

	return item, ttl_set
}

func (c *MemCache[T])__set_and_wait_if_need(key, value interface{}, cost int64, ttl time.Duration)bool{
	if c.add.Add(1) >= 4096 {
		c.c.Wait()
		c.add.Store(0)
	}
	return c.c.SetWithTTL(key, value, cost, ttl)
}

// __rerent_item_if_need will rerent the iterm by checking internal policy. 
func (c *MemCache[T])__rerent_item_if_need(key any, item any) *memCacheItem[T] {
	myItem := item.(*memCacheItem[T])
	
	if c.opts.AutoReRent && c.opts.MaxTTL > 0 {
		re_ttl := c.opts.MaxTTL

		if myItem.needExpiredAt > 0 {
			re_ttl = time.Until(time.UnixMilli(myItem.needExpiredAt))
			if re_ttl <= 0 {
				return myItem
			}
			if re_ttl > c.opts.MaxTTL {
				re_ttl = c.opts.MaxTTL
			}
		}

		cur_ttl, ok_ := c.c.GetTTL(key)
		if !ok_ {
			// this should not happen
			// internal ttl not get, that means the item expired in current time fo this operation is not atomic
			// here we return the item for it's got already 
			return myItem
		}
		// here to avoid rerent repeatedly by continues access
		if cur_ttl >= c.rerentSkipThresh{
			return myItem
		}
		// not need rerent if current_ttl > rerent_ttl
		if cur_ttl >= re_ttl {
			return myItem
		}

		c.__set_and_wait_if_need(key, item, int64(myItem.cost), re_ttl)
	}

	return myItem
}

func (c *MemCache[T])Set(key any, value T, ttl ...time.Duration) bool {
	item, ttl_set := c.__make_item(value, 0, ttl...)
	return c.__set_and_wait_if_need(key, item, 0, ttl_set)
}

func (c *MemCache[T])SetEx(key any, value T, cost int, ttl ...time.Duration) bool {
	item, ttl_set := c.__make_item(value, cost, ttl...)
	return c.__set_and_wait_if_need(key, item, int64(item.cost), ttl_set)
}

func (c *MemCache[T])Get(key any) (out T, ok bool) {
	item, _ := c.c.Get(key)
	if item == nil {
		return
	}

	return c.__rerent_item_if_need(key, item).Value, true
}

// return the left ttl of current key, if key not exist or expired then return 0, false
func (c *MemCache[T])GetTTL(key any) (time.Duration, bool) {
	return c.c.GetTTL(key)
}

func (c *MemCache[T])Del(key any) {
	c.c.Del(key)
}

func (c *MemCache[T])Dels(keys ...any) {
	for _, key := range keys{
		switch reflect.TypeOf(key).Kind() {
			case reflect.Slice, reflect.Array:
				s := reflect.ValueOf(key)
				len := s.Len()
				for i := 0; i < len; i++ {
					c.c.Del(s.Index(i))
				}
			default:
				c.c.Del(key)
	  }
	}
}

func (c *MemCache[T])Clear() {
	c.c.Clear()
	c.c.Wait()
}

func (c *MemCache[T])Close() {
	c.c.Wait()
	c.c.Close()
}

func (c *MemCache[T])Opt() *memCacheOpt[T] {
	return &memCacheOpt[T]{c: c}
}

type memCacheOpt[T any] struct {
	c *MemCache[T]
}

func (o *memCacheOpt[T])SetDefaultTTL(ttl time.Duration)*memCacheOpt[T]{
	if ttl > 0 {
		o.c.opts.TTL = ttl
	}
	return o
}

func (o *memCacheOpt[T])SetMaxTTL(ttl time.Duration)*memCacheOpt[T]{
	if ttl > 0 {
		o.c.opts.MaxTTL = ttl
		o.c.__cal_rerent_ttl()
	}
	return o
}

func (o *memCacheOpt[T])SetAutoReRent(on_off bool)*memCacheOpt[T]{
	o.c.opts.AutoReRent = on_off
	return o
}

func (o *memCacheOpt[T])ResetMaxcost(maxCost int64)*memCacheOpt[T]{
	o.c.opts.MaxCost = maxCost
	o.c.c.UpdateMaxCost(maxCost)
	return o
}