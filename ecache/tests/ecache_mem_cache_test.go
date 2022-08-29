package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

func TestMemCacheBaisc(t *testing.T){

  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)

	c.Set("k1", "v1")
	c.Set("k2", "v2")
	c.Wait()
  val, ok1 := c.Get("k1")
	ttl, ok2 := c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	c.Del("k1")
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, false, ok1)
	assert.Equal(t, false, ok2)

	c.Set("k1", "v1", time.Second)
	c.Set("k2", "v2", time.Second)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	c.Clear()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, false, ok1)
	assert.Equal(t, false, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, false, ok1)
	assert.Equal(t, false, ok2)
}

func TestMemCacheDfTTL(t *testing.T){
  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
		  TTL        : time.Second,
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)

	// not set ttl, using default
	c.Set("k1", "v1")
	c.Set("k2", "v2")
	c.Wait()
  val, ok1 := c.Get("k1")
	ttl, ok2 := c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	// set ttl, using the set
	c.Set("k1", "v1", 0)
	c.Set("k2", "v2", 0)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	c.Clear()
}

func TestMemCacheMaxTTL(t *testing.T){
  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
		  TTL        : time.Second,
			MaxTTL     : time.Second * 3,
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)
	defer c.Clear()

	// not set ttl, using default
	c.Set("k1", "v1")
	c.Set("k2", "v2")
	c.Wait()
  val, ok1 := c.Get("k1")
	ttl, ok2 := c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	// set ttl > MaxTTL, using MaxTTL
	c.Set("k1", "v1", time.Second * 10)
	c.Set("k2", "v2", time.Second * 10)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 3)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 3)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	// set ttl < MaxTTL, using the set
	c.Set("k1", "v1", time.Second * 2)
	c.Set("k2", "v2", time.Second * 2)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	// set ttl == 0, using the MaxTTL
	c.Set("k1", "v1", 0)
	c.Set("k2", "v2", 0)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 3)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 3)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	c.Clear()
  c = ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
		  TTL        : time.Second * 3,
			MaxTTL     : time.Second * 2,   // MaxTTL should take controll
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)

	// not set ttl, using MaxTTL because MaxTTL < TTL
	c.Set("k1", "v1")
	c.Set("k2", "v2")
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)

	// set ttl > MaxTTL, using MaxTTL
	c.Set("k1", "v1", time.Second * 10)
	c.Set("k2", "v2", time.Second * 10)
	c.Wait()
  val, ok1 = c.Get("k1")
	ttl, ok2 = c.GetTTL("k1")
	assert.Equal(t, "v1", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
  val, ok1 = c.Get("k2")
	ttl, ok2 = c.GetTTL("k2")
	assert.Equal(t, "v2", val)
	assert.Equal(t, true, ttl > 0)
	assert.Equal(t, true, ttl < time.Second * 2)
	assert.Equal(t, true, ok1)
	assert.Equal(t, true, ok2)
}

func TestMemCacheMaxCost(t *testing.T){
  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
			MaxCost    : 10,                   // 
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)
	defer c.Clear()

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111"},
		{Name:"name2", Tel:"22222222222"},
		{Name:"name3", Tel:"33333333333"},
		{Name:"name4", Tel:"44444444444"},
		{Name:"name5", Tel:"55555555555"},
		{Name:"name6", Tel:"66666666666"},
		{Name:"name7", Tel:"77777777777"},
		{Name:"name8", Tel:"88888888888"},
		{Name:"name9", Tel:"99999999999"},
		{Name:"name0", Tel:"00000000000"},
	}

	for _, i := range inputs {
		c.SetEx(i.Name, i.Tel, 1)
	}
	c.Wait()

	for _, i := range inputs {
		val, _ := c.Get(i.Name)
		assert.Equal(t, i.Tel, val)
	}
	c.SetEx("k1", "v1", 1)
	c.Wait()

	founds := 0
	for _, i := range inputs {
		_, ok := c.Get(i.Name)
		if ok == false {
			t.Logf("%s has been evicted ", i.Name )
		} else {
			founds++
		}
	}
	assert.Equal(t, founds, 9)
	v1, _ := c.Get("k1")
	assert.Equal(t, "v1", v1)

	fmt.Println(c.Metrics.String())
	assert.Equal(t, uint64(11), c.Metrics.KeysAdded())
	assert.Equal(t, uint64(1),  c.Metrics.KeysEvicted())

	c.Clear();

	// -------------------------------
	// test rejecting
	// 
	c = ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
			MaxCost    : 10,                   // 
			Statistics : true,
			IgnoreInternalCost: true,
		},
	)
	for _, i := range inputs {
		c.SetEx(i.Name, i.Tel, 1)
	}
	c.Wait()

	for _, i := range inputs {
		c.Get(i.Name)
	}
	c.SetEx("k1", "v1", 1)
	c.Wait()
	founds = 0
	for _, i := range inputs {
		_, ok := c.Get(i.Name)
		if ok == false {
			t.Logf("%s has been evicted ", i.Name )
		} else {
			founds++
		}
	}
	assert.Equal(t, founds, 9)
}

func TestMemCacheAutoRerent(t *testing.T){
  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
			MaxTTL     : time.Second * 2,
			Statistics : true,
			AutoReRent : true,
			IgnoreInternalCost: true,
		},
	)
	defer c.Clear()

	c.Set("k1", "v1", time.Second * 10)
	c.Set("k2", "v2", time.Second * 10)
	c.Set("k3", "v3", time.Second * 4)
	c.Set("k4", "v4")

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.Get("k2")   // here access 
		c.Get("k3")
		c.Get("k4")
	}
	c.Wait()

	v1, _ := c.Get("k1")
	t1, _ := c.GetTTL("k1")
	assert.Equal(t, "", v1)
	assert.Equal(t, true, t1 == 0)

	v2, _ := c.Get("k2")
  t2, _ := c.GetTTL("k2")
	assert.Equal(t, "v2", v2)
	assert.Equal(t, true, t2 > time.Second)
	assert.Equal(t, true, t2 < time.Second * 2)

	v3, _ := c.Get("k3")
	t3, _ := c.GetTTL("k3")
	assert.Equal(t, "v3", v3)
	assert.Equal(t, true, t3 > time.Second * 0)
	assert.Equal(t, true, t3 < time.Second * 1)

	v4, _ := c.Get("k4")
	t4, _ := c.GetTTL("k4")
	assert.Equal(t, "v4", v4)
	assert.Equal(t, true, t4 > time.Second * 1)
	assert.Equal(t, true, t4 < time.Second * 2)
}

func TestMemCacheOnEvict(t *testing.T){

	evicted := 0

  c := ecache.NewMemCache(
		ecache.MemCacheOpts[string]{
			MaxCost    : 10,
			Statistics : true,
			IgnoreInternalCost: true,
			OnEvict    : func(val string){
				//t.Logf("%s evicted\n", val)
				evicted += 1
			},
		},
	)
	defer c.Clear()

	for i := 0; i < 1000; i++ {
		kv := fmt.Sprintf("%d", i)
		c.SetEx(kv, kv, 1)
	}
	c.Wait()

	t.Logf("evicted %d\n", evicted)
	assert.Equal(t, 990, evicted)
}