package tests

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
	"github.com/ziyht/eden_go/etimer"
	"github.com/ziyht/eden_go/eutils/ptr"
)

//var nilVal []byte = nil
var val512 []byte

func TestAll(t *testing.T){
	for i := 0; i < 512; i++ {
		val512 = append(val512, byte(i % 10))
	}

	ExecTestForDsn(t, "badger:test_data/badger")
	ExecTestForDsn(t, "nutsdb:test_data/nutsdb")
}

func printCurAlloc(t *testing.T) {
  var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	
	t.Logf("cur alloc: %d", stat.Alloc / 1024 / 1024)
}

func TestBadger(t *testing.T){
	for i := 0; i < 512; i++ {
		val512 = append(val512, byte(i % 10))
	}

	etimer.AddInterval(nil, time.Second/10, func(j *etimer.Job)error {
		printCurAlloc(t)
		return nil
	})

	ExecInsert(t, "badger:test_data/badger", 5000000)
}

func TestNutsDB(t *testing.T){
	for i := 0; i < 512; i++ {
		val512 = append(val512, byte(i % 10))
	}

	etimer.AddInterval(nil, time.Second/10, func(j *etimer.Job)error {
		printCurAlloc(t)
		return nil
	})

	ExecInsert(t, "badger:test_data/nutsdb", 5000000)
}

func ExecTestForDsn(t *testing.T, dsn string){
	c, err := ecache.NewDBCache(ecache.DBCacheOpts{Dsn: dsn} )
	if !assert.Equal(t, nil, err){
		return
	}
	c.Truncate()
	c.Close()

	t.Logf("------------- %s ---------------", dsn)

	ExecInsert(t, dsn, 1000000)
  ExecSet512(t, dsn, 1000000)
	ExecSets512(t, dsn, 1000000)
}

func ExecInsert(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(ecache.DBCacheOpts{Dsn: dsn} )
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	start := time.Now()
	for _, k := range keys {
		r.Set(k, k)
	}
	t.Logf("set %d keys(v = k) for each: cost %s", cnt, time.Since(start))

	c.Truncate()
	c.Close()
}

func ExecSet512(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(ecache.DBCacheOpts{Dsn: dsn} )
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	start := time.Now()
	for _, k := range keys {
		r.Set(k, val512)
	}
	t.Logf("set %d keys(512val) for each: cost %s", cnt, time.Since(start))

	c.Truncate()
	c.Close()
}

func ExecSets512(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(ecache.DBCacheOpts{Dsn: dsn} )
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	start := time.Now()
	err = r.Sets(keys, val512)
	assert.Equal(t, nil, err)
	t.Logf("set %d keys(512val) by group: cost %s", cnt, time.Since(start))

	var vals []ecache.Val
	var val ecache.Val
	start = time.Now()
	for _, k := range keys {
		val, err = r.Get(k)
		assert.Equal(t, nil, err)
		vals = append(vals, val)
	}
	assert.Equal(t, cnt, len(vals))
	t.Logf("get %d keys(512val) for each: cost %s", len(vals), time.Since(start))
	for _, val := range vals {
		if !assert.Equal(t, val512, val.Bytes()){
			return
		}
	}

	start = time.Now()
	vals, err = r.Gets(keys)
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(vals))
	t.Logf("get %d keys(512val) by group: cost %s", len(vals), time.Since(start))

	start = time.Now()
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(keys))
	assert.Equal(t, cnt, len(vals))
	t.Logf("get %d keys(512val) by getall: cost %s", len(keys), time.Since(start))

	// c.Truncate()
	// c.Close()
}

func toBytesKey(key any)([]byte, error){
	switch k := key.(type) {
	case string: return ptr.StringToBytes(k), nil
	case []byte: return k, nil
	}
	return nil, fmt.Errorf("invalid key type, only support string and []byte")
}

func strToBytes(s string)([]byte){
	return []byte(s)
}

func bytesToBytes(b []byte)([]byte){
	return []byte(b)
}

func TestConvertKey(t *testing.T){

	cnt := 10000000

	skey := "1234567890123456789012345678901234567"
	bkey := []byte(skey)

	start := time.Now()
	for i := 0; i < cnt; i++ {
		toBytesKey(skey)
	}
	t.Logf("toBytesKey skey: %s", time.Since(start))

	start = time.Now()
	for i := 0; i < cnt; i++ {
		toBytesKey(bkey)
	}
	t.Logf("toBytesKey bkey: %s", time.Since(start))

	start = time.Now()
	for i := 0; i < cnt; i++ {
		strToBytes(skey)
	}
	t.Logf("strToBytes bkey: %s", time.Since(start))

	start = time.Now()
	for i := 0; i < cnt; i++ {
		bytesToBytes(bkey)
	}
	t.Logf("bytesToBytes bkey: %s", time.Since(start))

}