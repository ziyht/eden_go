package ecache

import (
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
)

type Metrics = ristretto.Metrics

type MemCache[V any] struct {
  c                *ristretto.Cache
	reRentSkipThresh time.Duration      // default = MaxTTL * 0.8
	reRentTTL        time.Duration      // default = MaxTTL
	opts             MemCacheOpts[V]
	add              atomic.Int32
	Metrics          *ristretto.Metrics
}

type MemCacheOpts[V any] struct {
	MaxCost            int64             // default:   10M, each item has a cost of memory, this value defines the max cost of current cache instance
	DfTTL              time.Duration     // default:     0, the default TTL of items when not passed ttl in params in Set functions, 0 means the item will never expired by TTL policy.
  MaxTTL             time.Duration     // default:     0, the max TTL of all items, the ttl will set to it if the passed ttl is greater than MaxTTL
	AutoReRent         bool              // default: false, the ttl of the item will be automatically increasing to MaxTTL(if is set) after access it
	OnDelete           func(any)         // default:   nil, this will be called whenever a value is removed from cache
	OnCost             func(any) int64   // default:   nil, this will be called to get real cost of item when the input cost is 0
	IgnoreInternalCost bool              // default: false, IgnoreInternalCost set to true indicates to the cache that the cost of internally storing the value should be ignored.
	CountersNum        int64             // default:  4096, CountersNum determines the number of counters (keys) to keep that hold access frequency information used by internal policy, not the count limit for items. It's generally a good idea to have more counters(10x) than the max cache capacity, as this will improve eviction accuracy and subsequent hit ratios.
	BufferItems        int64             // default:    64, ristretto: BufferItems is the size of the Get buffers. The best value we've found for this is 64.
	Statistics         bool              // default: false, do Statistics interval or not
}

func CostN(cost int64) func(any) int64  {
	return func(any) int64 { return cost }
}

func __cost_func_factory(v_ any) func(any) int64 {
	v := reflect.ValueOf(v_)
	switch v.Kind() {
	case reflect.Struct : return CostN(int64(unsafe.Sizeof(v.Elem())))
	case reflect.Ptr, reflect.Interface:
			p := (*[]byte)(unsafe.Pointer(v.Pointer()))
			if p == nil {
					return CostN(8)
			}
			return CostN(int64(unsafe.Sizeof(v.Elem())) + 8)
	case reflect.String, reflect.Array, reflect.Chan, reflect.Slice:
			return func(v_ any) int64 { v := reflect.ValueOf(v_); return int64(v.Len()) }

	case reflect.Bool, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Int:
			return CostN(int64(v.Type().Size()))

	default :
		return CostN(int64(unsafe.Sizeof(v)))
	}
}

func newMemCache[V any](opts MemCacheOpts[V]) *MemCache[V] {

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
		var v V
		opts.OnCost = __cost_func_factory(v)
	}

	c := &MemCache[V]{
		opts     : opts,
	}

	if opts.MaxTTL > 0 {
		if opts.DfTTL > opts.MaxTTL {
			opts.DfTTL = opts.MaxTTL
		}
		c.reRentSkipThresh = time.Duration(float64(opts.MaxTTL) * 0.8)
		c.reRentTTL = opts.MaxTTL
	}
	if opts.DfTTL > 0 {
		c.reRentSkipThresh = time.Duration(float64(opts.DfTTL) * 0.8)
		c.reRentTTL = opts.DfTTL
	}

	config := &ristretto.Config{
		NumCounters       : opts.CountersNum,
		MaxCost           : opts.MaxCost,
		BufferItems       : opts.BufferItems,
		IgnoreInternalCost: opts.IgnoreInternalCost,
		Cost              : opts.OnCost,
		Metrics           : opts.Statistics,
	}

	config.OnExit = opts.OnDelete

	c.c, _ = ristretto.NewCache(config)
	if c.c == nil {
		return nil
	}

	c.Metrics = c.c.Metrics

	return c
}

// the internal cache is not synchronized with Set funcs
// you should call Wait() after call Set funcs
func (c *MemCache[V])Wait() {
	c.c.Wait()
}

func (c *MemCache[V])__validate_ttl(ttl ...time.Duration) time.Duration {
	ttl_to_set := c.opts.DfTTL
	if len(ttl) > 0 {
		ttl_to_set = ttl[0]

		if ttl_to_set <= 0 {
			return time.Duration(0)
		}
	}

	// limited to max_ttl
	if c.opts.MaxTTL > 0 && (c.opts.MaxTTL < ttl_to_set || ttl_to_set == 0) {
		ttl_to_set = c.opts.MaxTTL
	}

	return ttl_to_set
}

func (c *MemCache[V])__set_and_wait_if_need(key, value interface{}, cost int64, ttl time.Duration)bool{
	if c.add.Add(1) >= 4096 {
		c.c.Wait()
		c.add.Store(0)
	}
	return c.c.SetWithTTL(key, value, cost, ttl)
}

// __rerent_item_if_need will rerent the iterm by checking internal policy. 
func (c *MemCache[V])__rerent_item_if_need(key any, val any) {
	if c.opts.AutoReRent && c.reRentTTL > 0 {
		ttl, ok := c.c.GetTTL(key)
		if !ok || ttl >= c.reRentSkipThresh{
			return
		}

		c.__set_and_wait_if_need(key, val, 0, c.reRentTTL)
	}
}

// Set stores the given value under the specified key in the cache.
// If a TTL is provided, it sets the time-to-live for the cache entry.
// If TTL is not provided, the default TTL for the cache is used.
// if first given TTL is 0, the value never expires by TTL policy, but it still can be removed by LRU policy.
// Returns true if the entry was successfully stored.
// this operation is not synchronized, you may not found the value immediately, we will wait synchronize internal in every 4096 calls, or you can call Wait() to synchronize manually
func (c *MemCache[V])Set(key any, val V, ttl ...time.Duration) bool {
	return c.__set_and_wait_if_need(key, val, 0, c.__validate_ttl(ttl...))
}

// SetSync is like Set, but it will block until the internal cache is synchronized with the set operation.
// It will return true if the entry was successfully stored.
func (c *MemCache[V])SetSync(key any, val V, ttl ...time.Duration) bool {
	defer c.c.Wait()
	return c.__set_and_wait_if_need(key, val, 0, c.__validate_ttl(ttl...))
}

func (c *MemCache[V])SetEx(key any, val V, cost int64, ttl ...time.Duration) bool {
	return c.__set_and_wait_if_need(key, val, cost, c.__validate_ttl(ttl...))
}

// Get will return the value which bind to the input key,
// if no value bind to the key, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[V])Get(key any) (val V, exist bool) {
	i, exist := c.c.Get(key)
	if exist {
		c.__rerent_item_if_need(key, i)
		return i.(V), exist
	}

	return
}

// GetByKeys will return the value which bind to the first key in the input keys,
// if no key bind any value, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[V]) GetByKeys(keys ...any) (val V, exist bool) {
	for _, key := range keys{
		i, exist := c.c.Get(key)
		if exist {
			c.__rerent_item_if_need(key, val)
			return i.(V), exist
		}
	}
	return
}

// Val will return the value which bind to the input key,
// if no value bind to the key, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[V]) Val(key any) (V) {
	i, exist := c.c.Get(key)
	if exist {
		c.__rerent_item_if_need(key, i)
		return i.(V)
	}

	var v V
	return v
}

// ValByKeys will return the value which bind to the first key in the input keys,
// if no key bind any value, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[V]) ValByKeys(keys ...any) (val V) {
	for _, key := range keys{
		i, exist := c.c.Get(key)
		if exist {
			c.__rerent_item_if_need(key, val)
			return i.(V)
		}
	}
	return
}

// return the left ttl of current key, if key not exist or expired then return 0, false
func (c *MemCache[V])GetTTL(key any) (time.Duration, bool) {
	return c.c.GetTTL(key)
}

func (c *MemCache[V])Del(key any) {
	c.c.Del(key)
}

func (c *MemCache[V])Dels(keys ...any) {
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

func (c *MemCache[V])Clear() {
	c.c.Clear()
	c.c.Wait()
}

func (c *MemCache[V])Close() {
	c.c.Wait()
	c.c.Close()
}

func (c *MemCache[V])Opt() *memCacheOpt[V] {
	return &memCacheOpt[V]{c: c}
}

type memCacheOpt[V any] struct {
	c *MemCache[V]
}

func (o *memCacheOpt[V])SetDefaultTTL(ttl time.Duration)*memCacheOpt[V]{
	if ttl > 0 {
		o.c.opts.DfTTL = ttl
	}
	return o
}

func (o *memCacheOpt[V])SetAutoReRent(on_off bool)*memCacheOpt[V]{
	o.c.opts.AutoReRent = on_off
	return o
}

func (o *memCacheOpt[V])SetMaxCost(maxCost int64)*memCacheOpt[V]{
	o.c.opts.MaxCost = maxCost
	o.c.c.UpdateMaxCost(maxCost)
	return o
}