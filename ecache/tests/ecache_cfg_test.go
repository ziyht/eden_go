package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ziyht/eden_go/ecache"
)

func TestCfgBasic(t *testing.T) {

	err := ecache.InitFromConfigFile("./cfgs/config1.yml")
	assert.NoError(t, err)
	err  = ecache.InitFromConfigFile("./cfgs/config2.yml")
	assert.NoError(t, err)

	dfc, e := ecache.GetDBCache()
	assert.NoError(t, e)
	assert.NotNil(t, dfc)

	c1, e1 := ecache.GetDBCache("cache1")
	c2, e2 := ecache.GetDBCache("cache2")
	assert.NoError(t, e1)
	assert.NoError(t, e2)
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)

	assert.NoError(t, dfc.Close())
	assert.NoError(t, c1.Close())
	assert.NoError(t, c2.Close())

	c3, e3 := ecache.GetDBCacheFromFile("./cfgs/config1.yml", "ecache.cache1")
	assert.NoError(t, e3)
	assert.NotNil(t, c3)
	assert.NoError(t, c3.Close())
}

func TestOptsBasic(t *testing.T){
	c1, e1 := ecache.NewDBCache(ecache.DBCacheOpts{Dir: "./test_data/cache5"})
	c2, e2 := ecache.NewDBCache(ecache.DBCacheOpts{Dir: "./test_data/cache6", Driver: ecache.BADGER})
	assert.NoError(t, e1)
	assert.NoError(t, e2)
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
}