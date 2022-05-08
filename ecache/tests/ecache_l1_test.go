package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)


func TestL1All(t *testing.T){
	ExecTestL1ForDsn(t, "badger:test_data/badger")
	ExecTestL1ForDsn(t, "nutsdb:test_data/nutsdb")
}

func ExecTestL1ForDsn(t *testing.T, dsn string){
  ExecTestL1_Basic(t, dsn)
	ExecTestL1_GetFromMem(t, dsn)
}

func ExecTestL1_Basic(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*2},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := c.GetL1(""  , time.Second*2)
	l11 := c.GetL1("l1", time.Second*2)
	l12 := c.GetL1("l2", time.Second*3)

	assert.Equal(t, true , l10 == c.GetL1("", time.Duration(0)))
	assert.Equal(t, true , l11 == c.GetL1("l1", time.Duration(0)))
	assert.Equal(t, true , l12 == c.GetL1("l2", time.Duration(0)))
	assert.Equal(t, false, l12 == c.GetL1("l1", time.Duration(0)))

	for _, i := range inputs {
		l10.SetItem(i.Name, i)
		l11.SetItem(i.Name, i)
		l12.SetItem(i.Name, i)
	}

	for _, i := range inputs {
		{
			ret, err := l10.GetItem(i.Name, i)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}

		{
			ret, err := l11.GetItem(i.Name, i)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}
		
		{
			ret, err := l12.GetItem(i.Name, i)
			assert.Equal(t, err, nil)
			get, ok := ret.(*myItem)
			assert.Equal(t, ok, true)
			assert.Equal(t, i.Equal(get), true)
		}
	}

	time.Sleep(time.Second*2)
	for _, i := range inputs {
		{
			ret, err := l10.GetItem(i.Name, i)
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
			ret, err := l11.GetItem(i.Name, i)
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
			ret, err := l12.GetItem(i.Name, i)
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

func ExecTestL1_GetFromMem(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	inputs := []*myItem{
		{Name:"name1", Tel:"11111111111", TTL_: time.Second*1},
		{Name:"name2", Tel:"22222222222", TTL_: time.Second*3},
		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
	}

	l10 := c.GetL1(""  , time.Second*1)
	l11 := c.GetL1("l1", time.Second*3)
	l12 := c.GetL1("l2", time.Second*2)

	for _, i := range inputs {
		l10.SetItem(i.Name, i)
		l11.SetItem(i.Name, i)
		l12.SetItem(i.Name, i)
	}

	time.Sleep(time.Second*2)

	for _, i := range inputs {
		{
			ret, _ := l10.GetItem(i.Name, i)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
				assert.Equal(t, get.UnMarshal_, true) // unmarshal
			}
		}

		{
			ret, err := l11.GetItem(i.Name, i)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
				assert.Equal(t, get.UnMarshal_, false) // not unmarshal, should return from mem directly
			}
		}
		
		{
			ret, _ := l12.GetItem(i.Name, i)
			if i.TTL() <= time.Second * 2 {
				assert.Equal(t, ret, nil)
			} else {
				assert.Equal(t, err, nil)
				get, ok := ret.(*myItem)
				assert.Equal(t, ok, true)
				assert.Equal(t, i.Equal(get), true)
				assert.Equal(t, get.UnMarshal_, true) // unmarshal
			}
		}
	}

	c.Truncate()
	c.Close()
}