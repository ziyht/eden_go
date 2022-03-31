package tests

import (
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

	//ExecInsert(t, dsn, 1000000)
  ExecSet512(t, dsn, 1000000)
	ExecSets512(t, dsn, 1000000)
}

func ExecInsert(t *testing.T, dsn string, cnt int){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	start := time.Now()
	for _, k := range keys {
		c.Set(k, k)
	}
	t.Logf("%s: insert %d keys: cost %s", dsn, cnt, time.Since(start))

	c.Truncate()
	c.Close()
}

func ExecSet512(t *testing.T, dsn string, cnt int){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	c.Truncate()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	v := make([]byte, 512)

	start := time.Now()
	for _, k := range keys {
		c.Set(k, v)
	}
	t.Logf("%s: insert %d keys: cost %s", dsn, cnt, time.Since(start))

	c.Close()
}

func ExecSets512(t *testing.T, dsn string, cnt int){
	c, err := ecache.GetCacheFromDsn(dsn)
	assert.Equal(t, nil, err)

	c.Truncate()

	keys := [][]byte{}
	for i := 0; i < cnt; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
	}

	v := make([]byte, 512)

	start := time.Now()
	err = c.Sets(keys, [][]byte{v})
	assert.Equal(t, nil, err)
	t.Logf("%s: insert %d keys: cost %s", dsn, cnt, time.Since(start))

	c.Close()
}