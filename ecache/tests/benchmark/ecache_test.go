package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
	"github.com/ziyht/eden_go/utils/ptr"
)

var nilVal []byte = nil

func TestAll(t *testing.T){
	//ExecTestForDsn(t, "badger:test_data/badger")
	ExecTestForDsn(t, "nutsdb:test_data/nutsdb")
}

func ExecTestForDsn(t *testing.T, dsn string){

	// ExecInsert(t, dsn, 1000000)
  ExecSet512(t, dsn, 1000000)
	ExecSets512(t, dsn, 1000000)
	ExecSets512Any(t, dsn, 1000000)
}

func ExecInsert(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(dsn)
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
	t.Logf("%s: insert %d keys: cost %s", dsn, cnt, time.Since(start))

	c.Truncate()
	c.Close()
}

func ExecSet512(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	c.Truncate()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	v := make([]byte, 512)

	start := time.Now()
	for _, k := range keys {
		r.Set(k, v)
	}
	t.Logf("%s: insert %d keys: cost %s", dsn, cnt, time.Since(start))

	c.Close()
}

func ExecSets512(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	c.Truncate()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	v := make([]byte, 512)

	start := time.Now()
	err = r.Sets(keys, [][]byte{v})
	assert.Equal(t, nil, err)
	t.Logf("%s: set by keys: %d keys: cost %s", dsn, cnt, time.Since(start))

	start = time.Now()
	vals, err := r.Gets(keys)
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(vals))
	t.Logf("%s: get by keys: %d keys: cost %s", dsn, len(vals), time.Since(start))

	start = time.Now()
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(keys))
	assert.Equal(t, cnt, len(vals))
	t.Logf("%s: get all    : %d keys: cost %s", dsn, len(keys), time.Since(start))

	c.Close()
}

func ExecSets512Any(t *testing.T, dsn string, cnt int){
	c, err := ecache.NewDBCache(dsn)
	assert.Equal(t, nil, err)
	r := c.DfRegion()

	c.Truncate()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	v := make([]byte, 512)

	start := time.Now()
	err = r.Sets(keys, [][]byte{v})
	assert.Equal(t, nil, err)
	t.Logf("%s: set by keys: %d keys: cost %s", dsn, cnt, time.Since(start))

	start = time.Now()
	vals, err := r.Gets(keys)
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(vals))
	t.Logf("%s: get by keys: %d keys: cost %s", dsn, len(vals), time.Since(start))

	start = time.Now()
	keys, vals, err = r.GetAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, cnt, len(keys))
	assert.Equal(t, cnt, len(vals))
	t.Logf("%s: get all    : %d keys: cost %s", dsn, len(keys), time.Since(start))

	c.Close()
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