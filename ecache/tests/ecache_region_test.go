package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

func TestRegion(t *testing.T){
	//ExecRegionTestForDsn(t, "badger:test_data/badger2")
	ExecRegionTestForDsn(t, "nutsdb:test_data/nutsdb")
}

func ExecRegionTestForDsn(t *testing.T, dsn string){
	ExecTestRegion_Basic(t, dsn)
	ExecTestRegion_Truncate(t, dsn)
	ExecTestRegion_SubRegion(t, dsn)
	ExecTestRegion_SubRegion2(t, dsn)
	ExecTestRegion_SubRegion3(t, dsn)
}


func ExecTestRegion_Basic(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()
	r1 := c.NewRegion("bb")  
	r2 := c.NewRegion("bbb") // 相同前缀的 region 数据 应该互不影响

	items := []any{
		&item{Key: "key1", Val: "val1", Ext: 0},
		&item{Key: "key2", Val: "val2", Ext: 0},
		&item{Key: "key3", Val: "val3", Ext: 0},
	}

	var keys[][]byte
	var keys1[][]byte
	var keys2[][]byte

	err = r.SetObjs(items, func(idx int, ite interface{})(k []byte, v any, du time.Duration){
		i := ite.(*item)
		keys = append(keys,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	// 这里写入到 r1，不应该影响其它 region 的数据
	err = r1.SetObjs(items, func(idx int, ite interface{})(k []byte, v any, du time.Duration){
		i := ite.(*item)
		keys1 = append(keys1,[]byte(i.Key))
		val, _ := json.Marshal(i)
		return []byte(i.Key), val, time.Duration(0)
	})
	assert.Equal(t, nil, err)

	// 这里写入到 r2，不应该影响其它 region 的数据
	err = r2.SetObjs(items, func(idx int, ite interface{})(k []byte, v any, du time.Duration){
		i := ite.(*item)
		keys2 = append(keys2,[]byte(i.Key))
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
	assert.Equal(t, len(gets), len(keys))
	for idx, item := range gets {
		assert.Equal(t, item, items[idx])
	}

	gets1 := make([]*item, 0)
	err = r1.DoForAll(func(idx int, k []byte, val ecache.Val)error{
		i := new(item)
		bin, err := val.GetBytes()  ; if err != nil { return err }
		err = json.Unmarshal(bin, i); if err != nil { return err }
		gets1 = append(gets1, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(gets1), len(keys1))
	for idx, item := range gets1 {
		assert.Equal(t, item, items[idx])
	}

	gets2 := make([]*item, 0)
	err = r1.DoForAll(func(idx int, k []byte, val ecache.Val)error{
		i := new(item)
		bin, err := val.GetBytes()  ; if err != nil { return err }
		err = json.Unmarshal(bin, i); if err != nil { return err }
		gets2 = append(gets2, i)
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, len(gets2), len(keys2))
	for idx, item := range gets2 {
		assert.Equal(t, item, items[idx])
	}

	c.Truncate()
	c.Close()
}

func ExecTestRegion_Truncate(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()
	r1 := c.NewRegion("bb")  
	r2 := c.NewRegion("bbb") // 相同前缀的 region 数据 应该互不影响

	keys_ := [][]byte{ []byte("key1"), []byte("key2"), []byte("key3")}
	vals_ := [][]byte{ []byte("val1"), []byte("val2"), []byte("val3")}

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	err = r1.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	err = r2.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 r1 的数据
	// =====================================
	err = r1.Truncate()
	assert.Equal(t, nil, err)

  // r 的数据不变
	keys, vals, err := r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys_), len(keys))
	assert.Equal(t, len(vals_), len(vals))

	// r1 的数据应该被清空
	keys, vals, err = r1.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	// r2 的数据应该不变
	keys, vals, err = r2.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys_), len(keys))
	assert.Equal(t, len(vals_), len(vals))

	// -------------------------------------
	// 清空 r2 的数据
	// =====================================
	err = r2.Truncate()
	assert.Equal(t, nil, err)

  // r 的数据应该不变
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys_), len(keys))
	assert.Equal(t, len(vals_), len(vals))

	// r1 的数据应该被清空
	keys, vals, err = r1.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	// r2 的数据应该被清空
	keys, vals, err = r2.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	// -------------------------------------
	// 清空 r 的数据
	// =====================================
	err = r.Truncate()
	assert.Equal(t, nil, err)

  // r 的数据应该被清空
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	// r1 的数据应该被清空
	keys, vals, err = r1.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	// r2 的数据应该被清空
	keys, vals, err = r2.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

	c.Truncate()
	c.Close()
}

func ExecTestRegion_SubRegion(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	r := c.DfRegion()

	r_s1 := r.SubRegion("s1")

	keys_ := [][]byte{ []byte("key1"), []byte("key2"), []byte("key3")}
	vals_ := [][]byte{ []byte("val1"), []byte("val2"), []byte("val3")}

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	err = r_s1.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 r 的数据
	// =====================================
	err = r.Truncate()
	assert.Equal(t, nil, err)

  // r 的数据应该被清空
	keys, vals, err := r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

  // r_s1 的数据应该不变
	keys, vals, err = r_s1.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys_), len(keys))
	assert.Equal(t, len(vals_), len(vals))

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	err = r_s1.Sets(keys_, vals_)
	assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 r_s1 的数据
	// =====================================
	err = r_s1.Truncate()
	assert.Equal(t, nil, err)

  // r_s1 的数据应该被清空
	keys, vals, err = r_s1.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(keys))
	assert.Equal(t, 0, len(vals))

  // r 的数据应该不变
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(keys_), len(keys))
	assert.Equal(t, len(vals_), len(vals))

	c.Truncate()
	c.Close()
}

func __subRegion2_check_helper(t *testing.T, r *ecache.Region, len_ int){
	keys, vals, err := r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, len_, len(keys))
	assert.Equal(t, len_, len(vals))
}

func ExecTestRegion_SubRegion2(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	r := c.DfRegion()
	r1 := c.NewRegion("r1")
	r2 := c.NewRegion("r2")
	r_s1 := r.SubRegion("s1")
	r_s2 := r.SubRegion("s2")
	r1_s1 := r1.SubRegion("s1")
	r1_s2 := r1.SubRegion("s2")
	r2_s1 := r2.SubRegion("s1")
	r2_s2 := r2.SubRegion("s2")


	keys_ := [][]byte{ []byte("key1"), []byte("key2"), []byte("key3")}
	vals_ := [][]byte{ []byte("val1"), []byte("val2"), []byte("val3")}

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r.Sets(keys_, vals_); assert.Equal(t, nil, err)
  err = r_s1.Sets(keys_, vals_);assert.Equal(t, nil, err)
  err = r_s2.Sets(keys_, vals_);assert.Equal(t, nil, err)
	err = r1.Sets(keys_, vals_); assert.Equal(t, nil, err)
  err = r1_s1.Sets(keys_, vals_);assert.Equal(t, nil, err)
  err = r1_s2.Sets(keys_, vals_);assert.Equal(t, nil, err)
	err = r2.Sets(keys_, vals_); assert.Equal(t, nil, err)
  err = r2_s1.Sets(keys_, vals_);assert.Equal(t, nil, err)
  err = r2_s2.Sets(keys_, vals_);assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 r 的数据
	// =====================================
	err = r.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, len(keys_))
	__subRegion2_check_helper(t, r_s2, len(keys_))
	__subRegion2_check_helper(t, r1_s1, len(keys_))
	__subRegion2_check_helper(t, r1_s2, len(keys_))
	__subRegion2_check_helper(t, r2_s1, len(keys_))
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r_s1 的数据
	// =====================================
	err = r_s1.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, len(keys_))
	__subRegion2_check_helper(t, r1_s1, len(keys_))
	__subRegion2_check_helper(t, r1_s2, len(keys_))
	__subRegion2_check_helper(t, r2_s1, len(keys_))
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r_s2 的数据
	// =====================================
	err = r_s2.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, 0)
	__subRegion2_check_helper(t, r1_s1, len(keys_))
	__subRegion2_check_helper(t, r1_s2, len(keys_))
	__subRegion2_check_helper(t, r2_s1, len(keys_))
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r1_s1 的数据
	// =====================================
	err = r1_s1.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, 0)
	__subRegion2_check_helper(t, r1_s1, 0)
	__subRegion2_check_helper(t, r1_s2, len(keys_))
	__subRegion2_check_helper(t, r2_s1, len(keys_))
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r1_s2 的数据
	// =====================================
	err = r1_s2.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, 0)
	__subRegion2_check_helper(t, r1_s1, 0)
	__subRegion2_check_helper(t, r1_s2, 0)
	__subRegion2_check_helper(t, r2_s1, len(keys_))
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r2_s1 的数据
	// =====================================
	err = r2_s1.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, 0)
	__subRegion2_check_helper(t, r1_s1, 0)
	__subRegion2_check_helper(t, r1_s2, 0)
	__subRegion2_check_helper(t, r2_s1, 0)
	__subRegion2_check_helper(t, r2_s2, len(keys_))

	// -------------------------------------
	// 清空 r2_s2 的数据
	// =====================================
	err = r2_s2.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r, 0)
	__subRegion2_check_helper(t, r_s1, 0)
	__subRegion2_check_helper(t, r_s2, 0)
	__subRegion2_check_helper(t, r1_s1, 0)
	__subRegion2_check_helper(t, r1_s2, 0)
	__subRegion2_check_helper(t, r2_s1, 0)
	__subRegion2_check_helper(t, r2_s2, 0)

	c.Truncate()
	c.Close()
}

func ExecTestRegion_SubRegion3(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)

	r1 := c.NewRegion("rr")
	r2 := c.NewRegion("rr", "rr")
	r1_s := r1.SubRegion("rr", "rr")
	r2_s := r1.SubRegion("rr")

	keys_ := [][]byte{ []byte("key1"), []byte("key2"), []byte("key3")}
	vals_ := [][]byte{ []byte("val1"), []byte("val2"), []byte("val3")}

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r1.Sets(keys_, vals_); assert.Equal(t, nil, err)
  err = r2.Sets(keys_, vals_); assert.Equal(t, nil, err)
	err = r1_s.Sets(keys_, vals_); assert.Equal(t, nil, err)
  err = r2_s.Sets(keys_, vals_); assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 r1_s 的数据
	// =====================================
	err = r1_s.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r1, len(keys_))
	__subRegion2_check_helper(t, r2, len(keys_))
	__subRegion2_check_helper(t, r1_s, 0)
	__subRegion2_check_helper(t, r2_s, len(keys_))

	// -------------------------------------
	// 写入数据
	// =====================================
	err = r1_s.Sets(keys_, vals_); assert.Equal(t, nil, err)

	// -------------------------------------
	// 清空 rs_s 的数据
	// =====================================
	err = r2_s.Truncate()
	assert.Equal(t, nil, err)
	__subRegion2_check_helper(t, r1, len(keys_))
	__subRegion2_check_helper(t, r2, len(keys_))
	__subRegion2_check_helper(t, r1_s, len(keys_))
	__subRegion2_check_helper(t, r2_s, 0)

	c.Truncate()
	c.Close()
}