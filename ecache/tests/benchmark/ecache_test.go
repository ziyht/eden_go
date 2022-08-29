package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

var nilVal []byte = nil

func TestAll(t *testing.T){
	ExecTestForDsn(t, "badger:test_data/badger")
	//ExecTestForDsn(t, "nutsdb:test_data/nutsdb")
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
		r.ASet(k, v)
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
	err = r.ASets(keys, [][]byte{v})
	assert.Equal(t, nil, err)
	t.Logf("%s: set by keys: %d keys: cost %s", dsn, cnt, time.Since(start))

	start = time.Now()
	vals, err := r.AGets(keys)
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
