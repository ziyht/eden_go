package ecache

import (
	"fmt"
	"sync"

	"github.com/ziyht/eden_go/elog"
)

var (
  dfDbCache *DBCache = nil
  dbCaches       = make(map[string]*DBCache)
  dbCacheCfgs    = make(map[string]*Cfg)
  dbCachesMu     = &sync.RWMutex{}
	dbCacheRootKey = "ecache"
	log            = elog.Log(elog.Opt().Filename("ecache"))
)

func initFromFile(path string) error {
	dbCachesMu.Lock()
	defer dbCachesMu.Unlock()

	cfgs, err := cfgsFromFile(path, dbCacheRootKey)
	if err != nil {
		return err
	}

	for name, cfg := range cfgs.Cfgs {
		fdcfg   := dbCacheCfgs[name]
		fdcache := dbCaches[name]

		if fdcache != nil {
			if fdcfg.Dsn != cfg.Dsn {
				return fmt.Errorf("the old cache's Dsn mismatch the new one set in file %s in %s.%s: %s != %s", path, dbCacheRootKey, name, cfg.Dsn, fdcfg.Dsn)
			}
			continue
		}

		if fdcfg != nil && fdcfg.Dsn != cfg.Dsn {
			log.Warnf("the old DSN '%s' from file '%s' overwriten by new DSN '%s' from file '%s' for key '%s'", fdcfg.Dsn, fdcfg.file, cfg.Dsn, cfg.file, name)
		}

		dbCacheCfgs[name] = cfg
	}
	return nil
}

func getDfDBcache() (*DBCache, error) {
	if dfDbCache == nil {
		tmp, err := NewDBCache(dfCfg().Dsn)
		if err != nil {
			return nil, err
		}
		dfDbCache = tmp
	}

	return dfDbCache, nil
}

func getDBCache(name string)(*DBCache, error){
	dbCachesMu.Lock()
	defer dbCachesMu.Unlock()

	cfg := dbCacheCfgs[name]
	if cfg == nil {
		return nil, fmt.Errorf("DBCache %s not found", name)
	}

	c := dbCaches[name]
	if c != nil {
		return c, nil
	}

	c, err := NewDBCache(cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("NewDBCache for '%s' failed: %s, dsn is: %s", name, err, cfg.Dsn)
	}

	dbCaches[name] = c
	return c, nil
}