package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

var nilVal []byte = nil

func TestBasic(t *testing.T){
	//ExecBasicTestForDsn(t, "badger:test_data/badger")
	ExecBasicTestForDsn(t, "nutsdb:test_data/nutsdb")
}

func ExecBasicTestForDsn(t *testing.T, dsn string){
	ExecTestBasic(t, dsn)
	ExecTestVal(t, dsn)
	ExecTestTTL(t, dsn)
	ExecTestClose(t, dsn)
	ExecTestSets(t, dsn)
	ExecTestIF_DoForKeys(t, dsn)
	ExecTestIF_DoForAll(t, dsn)

	// ExecTestBucketBasic(t, dsn)
	// ExecTestBucketIF_DoForKeys(t, dsn)
	// ExecTestBucketIF_DoForAll(t, dsn)

  // ExecTest_Item(t, dsn)
}

func ExecTestBasic(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	r := c.DfRegion()

	r.Set([]byte("key1"), []byte("value1"))
	r.Set([]byte("key2"), []byte("value2"), time.Millisecond)
	time.Sleep(time.Second)

	val1, _ := r.Get([]byte("key1"))
	val2, _ := r.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1.Bytes())
	assert.Equal(t, nilVal, val2.Bytes())

	val1, _ = r.Get([]byte("key1"))
	val2, _ = r.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1.Bytes())
	assert.Equal(t, nilVal, val2.Bytes())

	err = r.Dels([]byte("key1"), []byte("key2"))
	assert.Equal(t, nil, err)
	val1, _ = r.Get([]byte("key1"))
	val2, _ = r.Get([]byte("key2"))
	assert.Equal(t, nilVal, val1.Bytes())
	assert.Equal(t, nilVal, val2.Bytes())

	err = r.Set("key1", "value1", time.Second * 10)
	assert.Nil(t, err)
	time.Sleep(time.Second)
	v, expiresAt, err := r.GetEx("key1")
	assert.Nil(t, err)
	assert.Equal(t, "value1", v.Str())
	now := uint64(time.Now().Unix())
	diff := expiresAt - now
	assert.True(t, diff > 8)
	assert.True(t, diff < 10)
}

func ExecTestVal(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()
	
	r := c.DfRegion()
	
	r.Set("true" , true)
	r.Set("false", false)
	r.Set("int8" , int8(8))
	r.Set("int16", int16(16))
	r.Set("int32", int32(32))
	r.Set("int64", int64(64))
	r.Set("uint8" , uint8(8))
	r.Set("uint16", uint16(16))
	r.Set("uint32", uint32(32))
	r.Set("uint64", uint64(64))
	r.Set("f32"   , float32(32.0))
	r.Set("f64"   , float64(64.0))
	r.Set("bytes" , []byte("bytes"))
	r.Set("string", []byte("string"))
	r.Set("time"  , time.Now())
	r.Set("duration", time.Hour)

	b1,  err := r.Get("true");     assert.Equal(t, nil, err)
	b2,  err := r.Get("false");    assert.Equal(t, nil, err)
	i8,  err := r.Get("int8");     assert.Equal(t, nil, err)
	i16, err := r.Get("int16");    assert.Equal(t, nil, err)
	i32, err := r.Get("int32");    assert.Equal(t, nil, err)
	i64, err := r.Get("int64");    assert.Equal(t, nil, err)
	u8,  err := r.Get("uint8");    assert.Equal(t, nil, err)
	u16, err := r.Get("uint16");   assert.Equal(t, nil, err)
	u32, err := r.Get("uint32");   assert.Equal(t, nil, err)
	u64, err := r.Get("uint64");   assert.Equal(t, nil, err)
	f32, err := r.Get("f32");      assert.Equal(t, nil, err)
	f64, err := r.Get("f64");      assert.Equal(t, nil, err)
	bin, err := r.Get("bytes");    assert.Equal(t, nil, err)
	str, err := r.Get("string");   assert.Equal(t, nil, err)
	t1,  err := r.Get("time");     assert.Equal(t, nil, err)
	d,   err := r.Get("duration"); assert.Equal(t, nil, err)

  assert.Equal(t, true, b1.Bool())
	assert.Equal(t, false, b2.Bool())
	assert.Equal(t, int8(8), i8.I8())
	assert.Equal(t, int16(16), i16.I16())
	assert.Equal(t, int32(32), i32.I32())
	assert.Equal(t, int64(64), i64.I64())
	assert.Equal(t, uint8(8) ,  u8.U8())
	assert.Equal(t, uint16(16), u16.U16())
	assert.Equal(t, uint32(32), u32.U32())
	assert.Equal(t, uint64(64), u64.U64())
	assert.Equal(t, float32(32.0), f32.F32())
	assert.Equal(t, float64(64.0), f64.F64())
	assert.Equal(t, []byte("bytes"), bin.Bytes())
	assert.Equal(t, "string", str.Str())
	assert.Equal(t, true, time.Since(t1.Time()) < time.Second)
	assert.Equal(t, time.Hour, d.Duration())
}

func ExecTestTTL(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()

	r := c.DfRegion()

	r.Set([]byte("key1"), []byte("value1"))
	r.Set([]byte("key1"), []byte("value1"), time.Second)
	time.Sleep(time.Second)
	val1, _ := r.Get([]byte("key1"))
	assert.Equal(t, nilVal, val1.Bytes())

	r.Set([]byte("key1"), []byte("value1"), time.Second)
	r.Set([]byte("key1"), []byte("value1"))
	time.Sleep(time.Second)
	val1, _ = r.Get([]byte("key1"))
	assert.Equal(t, []byte("value1"), val1.Bytes())

	r.Set([]byte("key1"), []byte("value1"), time.Second)
	val1, _ = r.Get([]byte("key1"))
	assert.Equal(t, []byte("value1"), val1.Bytes())
}

func ExecTestClose(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	r := c.DfRegion()

	c.Truncate()
	r.Set([]byte("key1"), []byte("value1"))
	r.Set([]byte("key2"), []byte("value2"), time.Second)
	c.Close()

	c, err = ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	defer c.Close()
	defer c.Truncate()
	r = c.DfRegion()

	assert.Equal(t, nil, err)
	time.Sleep(time.Second)

	val1, _ := r.Get([]byte("key1"))
	val2, _ := r.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1.Bytes())
	assert.Equal(t, nilVal, val2.Bytes())

	val1, _ = r.Get([]byte("key1"))
	val2, _ = r.Get([]byte("key2"))
	assert.Equal(t, []byte("value1"), val1.Bytes())
	assert.Equal(t, nilVal, val2.Bytes())
}

func ExecTestSets(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	defer c.Close()
	defer c.Truncate()
	r := c.DfRegion()

	keys := [][]byte{}
	vals := [][]byte{}
	for i := 0; i < 100; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
		vals = append(vals, []byte(fmt.Sprintf("%d%d", i, i)))
	}

	err = r.Sets(keys, vals)
	assert.Equal(t, nil, err)

	gets, err := r.Gets(keys)
	assert.Equal(t, nil, err)
	for i := 0; i < 100; i++ {
		assert.Equal(t, vals[i], gets[i].Bytes())
	}
}

type item struct {
	Key    string
	Val    string
  Ext    time.Duration
}

func ExecTestIF_DoForKeys(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	r := c.DfRegion()

	objs := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	var keys[][]byte

	err = r.SetObjs(objs,  func(idx int, ite interface{})(k []byte, v any, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = r.DoForKeys(keys, func(idx int, k []byte, val ecache.Val)error{
		i := new(item)
		bin, err := val.GetBytes()  ; if err != nil { return err }
		err = json.Unmarshal(bin, i); if err != nil { return err }
		gets = append(gets, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(objs), len(keys))

	for idx, item := range gets {
		assert.Equal(t, item, objs[idx])
	}

	c.Truncate()
	c.Close()
}

func ExecTestIF_DoForAll(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: time.Second},
		&item{Key: "key2", Val: "val2", Ext: time.Second},
		&item{Key: "key3", Val: "val3", Ext: time.Second},
	}

	var keys[][]byte

	err = r.SetObjs(items, func(idx int, ite interface{})(k []byte, v any, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	gets := make([]*item, 0)
	err = r.DoForAll(func(idx int, k []byte, val ecache.Val)error{
		i := new(item)
		bin, err := val.GetBytes()  ; if err != nil { return err }
		err = json.Unmarshal(bin, i); if err != nil { return err }
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




// type testItem struct {
// 	User  string
// 	Phone string
// }

// func(i* testItem)New()ecache.Item{
// 	return new(testItem)
// }
// func(i* testItem)Key()[]byte{
// 	return []byte(i.User)
// }
// func(i* testItem)Marshal()([]byte, error){
// 	return json.Marshal(i)
// }
// func(i* testItem)Unmarshal(in []byte)(error){
// 	return json.Unmarshal(in, i)
// }
// func(i* testItem)TTL()(time.Duration){
// 	return time.Second
// }
// func(i* testItem)isEqual(in ecache.Item) bool {
// 	in_, ok := in.(*testItem)
// 	if ok == false {
// 		return false
// 	}
// 	if i.User != in_.User { return false }
// 	if i.Phone != in_.Phone { return false }
// 	return true
// }


// func ExecTestItem(t *testing.T, dsn string) {

// 	inputs := []*testItem{
// 		{User: "user1", Phone: "1234566"},
// 		{User: "user2", Phone: "16473323"},
// 		{User: "user3", Phone: "none"},
// 		{User: "user4", Phone: "76516334"},
// 	}

// 	c, err := ecache.GetCacheFromDsn(dsn)
// 	assert.Equal(t, nil, err)

// 	defer c.Close()
// 	defer c.Truncate()

// 	for _, item := range inputs {
// 		c.SetItem(item.Key(), item)
// 	}

// 	for _, item := range inputs {
// 		out_, err := c.GetItem(item.Key(), item)
// 		assert.Equal(t, err, nil)
// 		assert.Equal(t, item.isEqual(out_), true)
// 	}
// }


// func ExecTestBucketBasic(t *testing.T, dsn string){
// 	c, err := ecache.GetCacheFromDsn(dsn)
// 	assert.Equal(t, nil, err)

// 	defer c.Close()
// 	defer c.Truncate()

// 	c.BSet("b1", []byte("key1"), []byte("value1"))
// 	c.BSet("b2", []byte("key2"), []byte("value2"))
// 	assert.Equal(t, true, c.HaveBucket("b1"))
// 	assert.Equal(t, true, c.HaveBucket("b2"))
// 	assert.Equal(t, false, c.HaveBucket("b3"))

// 	val1, _ := c.BGet("b1", []byte("key1"))
// 	val2, _ := c.BGet("b2", []byte("key2"))
// 	assert.Equal(t, []byte("value1"), val1)
// 	assert.Equal(t, []byte("value2"), val2)

// 	val1, _ = c.BGet("b1", []byte("key1"))
// 	val2, _ = c.BGet("b2", []byte("key2"))
// 	assert.Equal(t, []byte("value1"), val1)
// 	assert.Equal(t, []byte("value2"), val2)

// 	c.BDel("b1", []byte("key1"))
// 	c.BDel("b2", []byte("key2"))
// 	val1, _ = c.BGet("b1", []byte("key1"))
// 	val2, _ = c.BGet("b2", []byte("key2"))
// 	assert.Equal(t, nilVal, val1)
// 	assert.Equal(t, nilVal, val2)
// }

// func ExecTestBucketIF_DoForKeys(t *testing.T, dsn string){
// 	c, err := ecache.GetCacheFromDsn(dsn)
// 	assert.Equal(t, nil, err)

// 	items := []any{
// 		&item{Key: "key1", Val: "val1", Ext: time.Second},
// 		&item{Key: "key2", Val: "val2", Ext: time.Second},
// 		&item{Key: "key3", Val: "val3", Ext: time.Second},
// 	}

// 	bucket1 := "bucket"
// 	bucket2 := "bucket2" // 具有相同前缀的 bucket 应该互不影响

// 	var keys1 [][]byte
// 	var keys2 [][]byte

// 	err = c.BSetIFs(bucket1, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
// 		i := ite.(*item)
// 		keys1 = append(keys1,[]byte(i.Key))
// 		val, _ := json.Marshal(i)
// 		return []byte(i.Key), val, time.Duration(0)
// 	})
// 	assert.Equal(t, nil, err)

// 	err = c.BSetIFs(bucket2, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
// 		i := ite.(*item)
// 		keys2 = append(keys2,[]byte(i.Key))
// 		val, _ := json.Marshal(i)
// 		return []byte(i.Key), val, time.Duration(0)
// 	})
// 	assert.Equal(t, nil, err)

// 	gets1 := make([]*item, 0)
// 	gets2 := make([]*item, 0)
// 	err = c.BDoForKeys(bucket1, keys1, func(idx int, k []byte, val []byte)error{
// 		i := new(item)
// 		err := json.Unmarshal(val, i)
// 		if err != nil {
// 			return err
// 		}
// 		gets1 = append(gets1, i)
// 		return nil
// 	})
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, len(gets1), len(keys1))

// 	err = c.BDoForKeys(bucket2, keys2, func(idx int, k []byte, val []byte)error{
// 		i := new(item)
// 		err := json.Unmarshal(val, i)
// 		if err != nil {
// 			return err
// 		}
// 		gets2 = append(gets2, i)
// 		return nil
// 	})
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, len(gets2), len(keys2))

// 	for idx, item := range gets1 {
// 		assert.Equal(t, item, items[idx])
// 	}
// 	for idx, item := range gets2 {
// 		assert.Equal(t, item, items[idx])
// 	}

// 	c.Truncate()
// 	c.Close()
// }

// func ExecTestBucketIF_DoForAll(t *testing.T, dsn string){
// 	c, err := ecache.GetCacheFromDsn(dsn)
// 	assert.Equal(t, nil, err)

// 	items := []any{
// 		&item{Key: "key1", Val: "val1", Ext: time.Second},
// 		&item{Key: "key2", Val: "val2", Ext: time.Second},
// 		&item{Key: "key3", Val: "val3", Ext: time.Second},
// 	}

// 	bucket1 := "bucket"
// 	bucket2 := "bucket2" // 具有相同前缀的 bucket 应该互不影响


// 	var keys1[][]byte
// 	var keys2[][]byte
// 	var keys3[][]byte

// 	err = c.BSetIFs(bucket1, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
// 		i := ite.(*item)
// 		keys1 = append(keys1,[]byte(i.Key))
// 		val, _ := json.Marshal(i)
// 		return []byte(i.Key), val, time.Duration(0)
// 	})
// 	assert.Equal(t, nil, err)

// 	err = c.BSetIFs(bucket2, items,  func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
// 		i := ite.(*item)
// 		keys2 = append(keys2,[]byte(i.Key))
// 		val, _ := json.Marshal(i)
// 		return []byte(i.Key), val, time.Duration(0)
// 	})
// 	assert.Equal(t, nil, err)

// 	// 写入到非 bucket ，不应对 bucket 数据产生影响
// 	err = c.SetIFs(items, func(idx int, ite interface{})(k []byte, v []byte, du time.Duration){
// 		i := ite.(*item)
// 		keys3 = append(keys3,[]byte(i.Key))
// 		val, _ := json.Marshal(i)
// 		return []byte(i.Key), val, time.Duration(0)
// 	})
// 	assert.Equal(t, nil, err)

// 	gets1 := make([]*item, 0)
// 	gets2 := make([]*item, 0)
// 	err = c.BDoForAll(bucket1, func(idx int, k []byte, val []byte)error{
// 		i := new(item)
// 		err := json.Unmarshal(val, i)
// 		if err != nil {
// 			return err
// 		}
// 		gets1 = append(gets1, i)
// 		return nil
// 	})
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, len(keys1), len(gets1))

// 	err = c.BDoForAll(bucket2, func(idx int, k []byte, val []byte)error{
// 		i := new(item)
// 		err := json.Unmarshal(val, i)
// 		if err != nil {
// 			return err
// 		}
// 		gets2 = append(gets2, i)
// 		return nil
// 	})
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, len(gets2), len(keys2))

// 	for idx, item := range gets1 {
// 		assert.Equal(t, item, items[idx])
// 	}
// 	for idx, item := range gets2 {
// 		assert.Equal(t, item, items[idx])
// 	}

// 	c.Truncate()
// 	c.Close()
// }



// func (i *myItem)Equal(i2 *myItem) bool {
// 	if i == i2 {
// 		return true
// 	}

// 	if i2 == nil {
// 		return false
// 	}

// 	if i.Name != i2.Name || i.Tel != i2.Tel || i.TTL_ != i2.TTL_{
// 		return false
// 	}
// 	return true
// }

// func ExecTest_Item(t *testing.T, dsn string){
// 	c, err := ecache.GetCacheFromDsn(dsn)
// 	assert.Equal(t, nil, err)

// 	inputs := []*myItem{
// 		{Name:"name1", Tel:"11111111111", TTL_: time.Second},
// 		{Name:"name2", Tel:"22222222222", TTL_: time.Second*2},
// 		{Name:"name3", Tel:"33333333333", TTL_: time.Second*3},
// 	}

// 	for _, i := range inputs {
// 		c.SetItem([]byte(i.Name), i, time.Second)
// 		c.BSetItem("b1", []byte(i.Name), i)
// 		c.BSetItem("b2", []byte(i.Name), i, time.Second)
// 	}

// 	for _, i := range inputs {
// 		{
// 			ret, err := c.GetItem([]byte(i.Name), i)
// 			assert.Equal(t, err, nil)
// 			get, ok := ret.(*myItem)
// 			assert.Equal(t, ok, true)
// 			assert.Equal(t, i.Equal(get), true)
// 		}

// 		{
// 			ret, err := c.BGetItem("b1", []byte(i.Name), i)
// 			assert.Equal(t, err, nil)
// 			get, ok := ret.(*myItem)
// 			assert.Equal(t, ok, true)
// 			assert.Equal(t, i.Equal(get), true)
// 		}
		
// 		{
// 			ret, err := c.BGetItem("b2", []byte(i.Name), i)
// 			assert.Equal(t, err, nil)
// 			get, ok := ret.(*myItem)
// 			assert.Equal(t, ok, true)
// 			assert.Equal(t, i.Equal(get), true)
// 		}
// 	}

// 	time.Sleep(time.Second)

// 	for _, i := range inputs {
// 		{
// 			ret, err := c.BGetItem("b1", []byte(i.Name), i)
// 			if i.TTL() <= time.Second {
// 				assert.Equal(t, ret, nil)
// 			}else {
// 				assert.Equal(t, err, nil)
// 				get, ok := ret.(*myItem)
// 				assert.Equal(t, ok, true)
// 				assert.Equal(t, i.Equal(get), true)
// 			}
// 		}
// 	}

// 	c.Truncate()
// 	c.Close()
// }