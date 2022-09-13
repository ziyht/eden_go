package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

type myItem struct {
	Name       string
	Tel        string
	TTL_       time.Duration
	UnMarshal_ bool
}

func newMyItem() ecache.Item {
	return &myItem{}
}

func newMyItem2() *myItem {
	return &myItem{}
}

func (i *myItem)Marshal()([]byte, error){
	return json.Marshal(i)
}

func (i *myItem)Unmarshal(in []byte)(error){
	err := json.Unmarshal(in, i)
	i.UnMarshal_ = true
	return err
}

func (i *myItem)TTL() time.Duration {
	return i.TTL_
}

func (i *myItem)Equal(i2 *myItem) bool {
	if i == i2 {
		return true
	}

	if i2 == nil {
		return false
	}

	if i.Name != i2.Name || i.Tel != i2.Tel || i.TTL_ != i2.TTL_{
		return false
	}
	return true
}

func TestItemRegionAll(t *testing.T){
	//ExecTestItemRegionDsn(t, "badger:test_data/badger")
	ExecTestItemRegionDsn(t, "nutsdb:test_data/nutsdb")
}

func ExecTestItemRegionDsn(t *testing.T, dsn string){
  ExecTestL1_Basic(t, dsn)
	ExecTestL1_MakeTypedRegion(t, dsn)
	ExecTestL1_GetFromMem(t, dsn)
	ExecTestL1_GetFromMemAfterReopen(t, dsn)
}

func ExecTestL1_Basic(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*2},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := c.NewItemRegion("")  ; l10.SetDefaultTTL(time.Second*2)
	l11 := c.NewItemRegion("l1"); l11.SetDefaultTTL(time.Second*2)
	l12 := c.NewItemRegion("l2"); l12.SetDefaultTTL(time.Second*3)

	for _, i := range inputs {
	  l10.Set(i.Name, i, i.TTL_)
		l11.Set(i.Name, i, i.TTL_)
		l12.Set(i.Name, i, i.TTL_)
	}

	for _, i := range inputs {
		{
			ret, err := l10.Get(i.Name, newMyItem)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}

		{
			ret, err := l11.Get(i.Name, newMyItem)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}
		
		{
			ret, err := l12.Get(i.Name, newMyItem)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}
	}

	time.Sleep(time.Second*2)
	for _, i := range inputs {
		{
			ret, err := l10.Get(i.Name, newMyItem)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
			}
		}

		{
			ret, err := l11.Get(i.Name, newMyItem)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
			}
		}
		
		{
			ret, err := l12.Get(i.Name, newMyItem)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
			}
		}
	}

	c.Truncate()
	c.Close()
}

func ExecTestL1_MakeTypedRegion(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*2},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := ecache.NewTypedItemRegion[*myItem](c.NewRegion(""))  ; l10.SetDefaultTTL(time.Second*2)
	l11 := ecache.NewTypedItemRegion[*myItem](c.NewRegion("l1")); l11.SetDefaultTTL(time.Second*2)
	l12 := ecache.NewTypedItemRegion[*myItem](c.NewRegion("l2")); l12.SetDefaultTTL(time.Second*3)

	for _, i := range inputs {
	  l10.Set(i.Name, i, i.TTL_)
		l11.Set(i.Name, i, i.TTL_)
		l12.Set(i.Name, i, i.TTL_)
	}

	for _, i := range inputs {
		{
			get, err := l10.Get([]byte(i.Name), newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}

		{
			get, err := l11.Get([]byte(i.Name), newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}
		
		{
			get, err := l12.Get([]byte(i.Name), newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}
	}

	var nilItem *myItem
	time.Sleep(time.Second*2)
	for _, i := range inputs {
		{
			get, err := l10.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}

		{
			get, err := l11.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}
		
		{
			get, err := l12.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}
	}

	c.Truncate()
	c.Close()
}

func ExecTestL1_GetFromMem(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second*1},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*3},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := ecache.NewTypedItemRegion[*myItem](c.NewRegion(""))  ; l10.SetDefaultTTL(time.Second*2)
	l11 := ecache.NewTypedItemRegion[*myItem](c.NewRegion("l1")); l11.SetDefaultTTL(time.Second*2)
	l12 := ecache.NewTypedItemRegion[*myItem](c.NewRegion("l2")); l12.SetDefaultTTL(time.Second*3)

	l10.EnableMemCache(10, time.Second * 3)
	l11.EnableMemCache(10, time.Second * 3)
	l12.EnableMemCache(10, time.Second * 3)

	for _, i := range inputs {
	  l10.Set(i.Name, i, i.TTL_)
		l11.Set(i.Name, i, i.TTL_)
		l12.Set(i.Name, i, i.TTL_)
	}

	for _, i := range inputs {
		{
			get, err := l10.Get(i.Name, newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}

		{
			get, err := l11.Get([]byte(i.Name), newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}
		
		{
			get, err := l12.Get([]byte(i.Name), newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}
	}

	var nilItem *myItem
	time.Sleep(time.Second*2)
	for _, i := range inputs {
		{
			get, err := l10.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}

		{
			get, err := l11.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}
		
		{
			get, err := l12.Get([]byte(i.Name), newMyItem2)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, get, nilItem)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, i.Equal(get), true)
			}
		}
	}

	t.Logf("l10: %s\n", l10.Metrics.String())
	t.Logf("l11: %s\n", l11.Metrics.String())
	t.Logf("l12: %s\n", l12.Metrics.String())

	c.Truncate()
	c.Close()
}

func ExecTestL1_GetFromMemAfterReopen(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second*1},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*3},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := ecache.NewTypedItemRegion[*myItem](c.NewRegion(""))  ; l10.SetDefaultTTL(time.Second*2)

	for _, i := range inputs {
	  l10.Set(i.Name, i)
	}

	err = c.Close()
	assert.Equal(t, nil, err)
	c, err = ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	l10 = ecache.NewTypedItemRegion[*myItem](c.NewRegion(""))  ; l10.SetDefaultTTL(time.Second*2)
	l10.EnableMemCache(10, time.Second * 2)
	for _, i := range inputs {
		get, err := l10.Get(i.Name, newMyItem2)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, i.Equal(get))
	}
	time.Sleep(time.Second / 10)
	for i := 0; i < 100; i++ {
		for _, i := range inputs {
			get, err := l10.Get(i.Name, newMyItem2)
			assert.Equal(t, err, nil)
			assert.Equal(t, i.Equal(get), true)
		}
	}
	assert.Equal(t, uint64(300), l10.Metrics.Hits())
	t.Logf("l10: %s\n", l10.Metrics.String())
}