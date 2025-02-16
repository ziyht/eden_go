package ecache

import (
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
)

type Key = z.Key
type Metrics = ristretto.Metrics

type MemCache[K Key, V any] struct {
  c                *ristretto.Cache[K, V]
	reRentSkipThresh time.Duration      // default = MaxTTL * 0.8
	reRentTTL        time.Duration      // default = MaxTTL
	opts             MemCacheOpts[V]
	add              atomic.Int32
	Metrics 				 *Metrics
}

type MemCacheOpts[T any] struct {
	MaxCost            int64             // default:   10M, each item has a cost of memory, this value defines the max cost of current cache instance
	DfTTL              time.Duration     // default:     0, the default TTL of items when not passed ttl in params in Set functions, 0 means the item will never expired by TTL policy, it will be limited by MaxTTL.
  MaxTTL             time.Duration     // default:     0, the max TTL of all items, the ttl will set to it if the passed ttl is greater than MaxTTL, 0 means no limit
	AutoReRent         bool              // default: false, the ttl of the item will be automatically increasing to MaxTTL(if is set) or DfTTL after access it
	OnDelete           func(T)           // default:   nil, this will be called whenever a value is removed from cache
	OnCost             func(T) int64     // default:   nil, this will be called to get real cost of item when the input cost is 0
	IgnoreInternalCost bool              // default: false, IgnoreInternalCost set to true indicates to the cache that the cost of internally storing the value should be ignored.
	CountersNum        int64             // default:  4096, CountersNum determines the number of counters (keys) to keep that hold access frequency information used by internal policy, not the count limit for items. It's generally a good idea to have more counters(10x) than the max cache capacity, as this will improve eviction accuracy and subsequent hit ratios.
	BufferItems        int64             // default:    64, ristretto: BufferItems is the size of the Get buffers. The best value we've found for this is 64.
	Statistics         bool              // default: false, do Statistics interval or not, 
}

func CostN[ T any](cost int64) func(T) int64  {
	return func(T) int64 { return cost }
}

func __cost_func_factory[T any](v_ T) func(T) int64 {
	v := reflect.ValueOf(v_)
	switch v.Kind() {
	case reflect.Struct : return CostN[T](int64(unsafe.Sizeof(v.Elem())))
	case reflect.Ptr, reflect.Interface:
			p := (*[]byte)(unsafe.Pointer(v.Pointer()))
			if p == nil {
					return CostN[T](8)
			}
			return CostN[T](int64(unsafe.Sizeof(v.Elem())) + 8)
	case reflect.String, reflect.Array, reflect.Chan, reflect.Slice:
			return func(v_ T) int64 { v := reflect.ValueOf(v_); return int64(v.Len()) }

	case reflect.Bool, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Int:
			return CostN[T](int64(v.Type().Size()))

	default :
		return CostN[T](int64(unsafe.Sizeof(v)))
	}
}
  
func newMemCache[K Key, V any](opts MemCacheOpts[V]) *MemCache[K, V] {

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

	c := &MemCache[K, V]{
		opts : opts,
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

	config := &ristretto.Config[K, V]{
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
// you should call Wait() after call Set funcs if you need to synchronize
func (c *MemCache[K, V])Wait() {
	c.c.Wait()
}

func (c *MemCache[K, V])__validate_ttl(ttl ...time.Duration) time.Duration {
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

func (c *MemCache[K, V])__wait_on_changes() {
	if c.add.Add(1) >= 4096 {
		c.c.Wait()
		c.add.Store(0)
	}
}

func (c *MemCache[K, V])__set_and_wait_if_need(key K, val V, cost int64, ttl time.Duration)bool{
	c.__wait_on_changes()
	return c.c.SetWithTTL(key, val, cost, ttl)
}

// __rerent_item_if_need will rerent the iterm by checking internal policy. 
func (c *MemCache[K, V])__rerent_item_if_need(key K, val V) {
	if c.opts.AutoReRent && c.reRentTTL > 0 {
		leftTTL, ok := c.c.GetTTL(key)
		if !ok || leftTTL >= c.reRentSkipThresh{
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
func (c *MemCache[K, V])Set(key K, val V, ttl ...time.Duration) bool {
	return c.__set_and_wait_if_need(key, val, 0, c.__validate_ttl(ttl...))
}

// SetSync is like Set, but it will block until the internal cache is synchronized with the set operation.
// It will return true if the entry was successfully stored.
func (c *MemCache[K, V])SetSync(key K, val V, ttl ...time.Duration) bool {
	defer c.c.Wait()
	return c.__set_and_wait_if_need(key, val, 0, c.__validate_ttl(ttl...))
}

func (c *MemCache[K, V])SetEx(key K, val V, cost int64, ttl ...time.Duration) bool {
	return c.__set_and_wait_if_need(key, val, cost, c.__validate_ttl(ttl...))
}

// Get will return the value which bind to the input key,
// if no value bind to the key, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[K, V])Get(key K) (val V, exist bool) {
	val, exist = c.c.Get(key)
	if exist {
		c.__rerent_item_if_need(key, val)
	}
	return
}

// GetByKeys will return the value which bind to the first key in the input keys, 
// if no key bind any value, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[K, V]) GetByKeys(keys ...K) (val V, exist bool) {
	for _, key := range keys{
		val, exist = c.c.Get(key)
		if exist {
			c.__rerent_item_if_need(key, val)
			return
		}
	}
	return
}

// Val will return the value which bind to the input key, 
// if no value bind to the key, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[K, V]) Val(key K) (V) {
	val, exist := c.c.Get(key)
	if exist {
		c.__rerent_item_if_need(key, val)
		return val
	}

	var v V
	return v
}

// ValByKeys will return the value which bind to the first key (which in cache) in the input keys,
// if no key bind any value, it will return the zero value of V.
// It will rerent the item if the config AutoReRent is true and the item is not expired.
func (c *MemCache[K, V]) ValByKeys(keys ...K) (val V) {
	var exist bool
	for _, key := range keys{
		val, exist = c.c.Get(key)
		if exist {
			c.__rerent_item_if_need(key, val)
			return
		}
	}
	return
}

// return the left ttl of current key, if key not exist or expired then return 0, false
func (c *MemCache[K, V])GetTTL(key K) (time.Duration, bool) {
	return c.c.GetTTL(key)
}

func (c *MemCache[K, V])Del(key K) {
	c.c.Del(key)
}

func (c *MemCache[K, V])Dels(keys ...K) {
	for _, key := range keys{
		c.c.Del(key)
	}
}

func (c *MemCache[K, V])Clear() {
	c.c.Clear()
	c.c.Wait()
}

func (c *MemCache[K, V])Close() {
	c.c.Wait()
	c.c.Close()
}

func (c *MemCache[K, V])Opt() MemCacheOpts[V] {
	return c.opts
}

func (c *MemCache[K, V])SetMaxCost(maxCost int64){
	c.c.UpdateMaxCost(maxCost)
}