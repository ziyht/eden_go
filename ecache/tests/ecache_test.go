package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
	_ "github.com/ziyht/eden_go/ecache/driver/badgerdb"
)

var nilVal []byte = nil

var dsn = "badger:test_data/badger"

func TestAll(t *testing.T){
	ExecTestForDsn(t, "badger:test_data/badger")
}

func ExecTestForDsn(t *testing.T, dsn string){
	ExecTestBasic(t, dsn)
	ExecTestClose(t, dsn)
	ExecTestIF_DoForKeys(t, dsn)
	ExecTestIF_DoForAll(t, dsn)

	ExecTestBucketBasic(t, dsn)
	ExecTestBucketIF_DoForKeys(t)
	ExecTestBucketIF_DoForAll(t)
}

func ExecTestBasic(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	c.Set([]byte("key1"), []byte("value1"))
	c.Set([]byte("key2"), []byte("value2"), time.Second)
	time.Sleep(time.Second)

	val1, _ := c.Get([]byte("key1"))
	val2, _ := c.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, nilVal, val2)

	val1, _ = c.Get([]byte("key1"))
	val2, _ = c.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, nilVal, val2)

	c.Dels([]byte("key1"), []byte("key2"))
	val1, _ = c.Get([]byte("key1"))
	val2, _ = c.Get([]byte("key2"))
	assert.Equal(t, nilVal, val1)
	assert.Equal(t, nilVal, val2)
}

func ExecTestClose(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	c.Truncate()
	c.Set([]byte("key1"), []byte("value1"))
	c.Set([]byte("key2"), []byte("value2"), time.Second)
	c.Close()

	c, err = ecache.GetCacheFromDsn(dsn)
	defer c.Close()
	defer c.Truncate()

	assert.Equal(t, nil, err)
	time.Sleep(time.Second)

	val1, _ := c.Get([]byte("key1"))
	val2, _ := c.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, nilVal, val2)

	val1, _ = c.Get([]byte("key1"))
	val2, _ = c.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, nilVal, val2)
}

type item struct {
	Key    string
	Val    string
  Ext    time.Duration
}

func ExecTestIF_DoForKeys(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	var keys[][]byte

	err = c.SetIFs(items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = c.DoForKeys(keys, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets = append(gets, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(items), len(keys))

	for idx, item := range gets {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}

func ExecTestIF_DoForAll(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	var keys[][]byte
	var keys2[][]byte

	err = c.SetIFs(items, func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	// 这里写入到 bucket，不应该影响非 bucket 数据
	err = c.BSetIFs("b", items, func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys2 = append(keys2,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = c.DoForAll(func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets = append(gets, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(items), len(keys))
	for idx, item := range gets {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}

func ExecTestBucketBasic(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	c.BSet("b1", []byte("key1"), []byte("value1"))
	c.BSet("b2", []byte("key2"), []byte("value2"))
	assert.Equal(t, true, c.HaveBucket("b1"))
	assert.Equal(t, true, c.HaveBucket("b2"))
	assert.Equal(t, false, c.HaveBucket("b3"))

	val1, _ := c.BGet("b1", []byte("key1"))
	val2, _ := c.BGet("b2", []byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, []byte("value2"), val2)

	val1, _ = c.BGet("b1", []byte("key1"))
	val2, _ = c.BGet("b2", []byte("key2"))
	assert.Equal(t, []byte("value1"), val1)
	assert.Equal(t, []byte("value2"), val2)

	c.BDel("b1", []byte("key1"))
	c.BDel("b2", []byte("key2"))
	val1, _ = c.BGet("b1", []byte("key1"))
	val2, _ = c.BGet("b2", []byte("key2"))
	assert.Equal(t, nilVal, val1)
	assert.Equal(t, nilVal, val2)
}

func ExecTestBucketIF_DoForKeys(t *testing.T){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	bucket := "bucket"

	var keys[][]byte

	err = c.BSetIFs(bucket, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = c.BDoForKeys(bucket, keys, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets = append(gets, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(items), len(keys))

	for idx, item := range gets {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}

func ExecTestBucketIF_DoForAll(t *testing.T){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	bucket := "bucket"

	var keys[][]byte
	var keys2[][]byte

	err = c.BSetIFs(bucket, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	// 写入到非 bucket ，不应对 bucket 数据产生影响
	err = c.SetIFs(items, func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys2 = append(keys2,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = c.BDoForAll(bucket, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets = append(gets, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(items), len(keys))

	for idx, item := range gets {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}