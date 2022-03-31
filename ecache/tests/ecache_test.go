package tests

import (
	"encoding/json"
	"fmt"
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
	ExecTestTTL(t, dsn)
	ExecTestClose(t, dsn)
	ExecTestSets(t, dsn)
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
	c.Set([]byte("key2"), []byte("value2"), time.Millisecond)
	time.Sleep(time.Millisecond * 10)

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

func ExecTestTTL(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	c.Set([]byte("key1"), []byte("value1"))
	c.Set([]byte("key1"), []byte("value1"), time.Second)
	time.Sleep(time.Second)
	val1, _ := c.Get([]byte("key1"))
	assert.Equal(t, nilVal, val1)

	c.Set([]byte("key1"), []byte("value1"), time.Second)
	c.Set([]byte("key1"), []byte("value1"))
	time.Sleep(time.Second)
	val1, _ = c.Get([]byte("key1"))
	assert.Equal(t, []byte("value1"), val1)

	c.Set([]byte("key1"), []byte("value1"), time.Second)
	val1, _ = c.Get([]byte("key1"))
	assert.Equal(t, []byte("value1"), val1)
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

func ExecTestSets(t *testing.T, dsn string){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	keys := [][]byte{}
	vals := [][]byte{}
	for i := 0; i < 100; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
		vals = append(vals, []byte(fmt.Sprintf("%d%d", i, i)))
	}

	err = c.Sets(keys, vals)

	gets, err := c.Gets(keys...)
	assert.Equal(t, nil, err)
	for i := 0; i < 100; i++ {
		assert.Equal(t, gets[i], vals[i])
	}
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

	bucket1 := "bucket"
	bucket2 := "bucket2" // 具有相同前缀的 bucket 应该互不影响

	var keys1 [][]byte
	var keys2 [][]byte

	err = c.BSetIFs(bucket1, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys1 = append(keys1,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	err = c.BSetIFs(bucket2, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys2 = append(keys2,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets1 := make([]*item, 0)
	gets2 := make([]*item, 0)
	err = c.BDoForKeys(bucket1, keys1, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets1 = append(gets1, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(gets1), len(keys1))

	err = c.BDoForKeys(bucket2, keys2, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets2 = append(gets2, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(gets2), len(keys2))

	for idx, item := range gets1 {
		assert.Equal(t, item, items[idx])
	}
	for idx, item := range gets2 {
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

	bucket1 := "bucket"
	bucket2 := "bucket2" // 具有相同前缀的 bucket 应该互不影响


	var keys1[][]byte
	var keys2[][]byte
	var keys3[][]byte

	err = c.BSetIFs(bucket1, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys1 = append(keys1,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	err = c.BSetIFs(bucket2, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys2 = append(keys2,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	// 写入到非 bucket ，不应对 bucket 数据产生影响
	err = c.SetIFs(items, func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
		i := ite.(*item)
		keys3 = append(keys3,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets1 := make([]*item, 0)
	gets2 := make([]*item, 0)
	err = c.BDoForAll(bucket1, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets1 = append(gets1, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys1), len(gets1))

	err = c.BDoForAll(bucket2, func(idx int, k []byte, val []byte)error{
		i := new(item)
		err := json.Unmarshal(val, i)
		if err != nil {
			return err
		}
		gets2 = append(gets2, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(gets2), len(keys2))

	for idx, item := range gets1 {
		assert.Equal(t, item, items[idx])
	}
	for idx, item := range gets2 {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}