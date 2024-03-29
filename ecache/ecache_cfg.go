package ecache

import (
	"fmt"
	"path/filepath"

	"github.com/ziyht/eden_go/ecfg"
)

type Cfg struct {
	Dsn     string   // if set, the other proporty will take no effect, format: <dirvername>:<dir>[?<arg1=val1>[&<arg2=val2>]...]
	Driver  string   // if not set will using nutsdb in default
	Dir     string   
  Params  map[string][]string
	cfgFile string
}

type Cfgs struct {
	Cfgs map[string]*Cfg
}

func dfCfg() *Cfg {
	return &Cfg{Dsn: "nutsdb:./cache/ecache/df/nutsdb"}
}

func cfgsFromFile(path string, key string)(Cfgs, error){
	var out Cfgs

	out.Cfgs = map[string]*Cfg{}

	path, err := filepath.Abs(path)
	if err != nil {
		return out, err
	}

	err = ecfg.ParsingFromCfgFile(path, key, &out.Cfgs)
	if err != nil {
		return out, fmt.Errorf("failed to parse %s from path %s: %s", key, path, err)
	}

	for _, cfg := range out.Cfgs{
		if cfg.Dsn == "" {
			cfg.Dsn = GenDsn(cfg.Driver, cfg.Dir, cfg.Params)
		}

		cfg.cfgFile = path
	}

	return out, nil
}

func cfgsFromFileKey(path string, key string)(Cfg, error){
	var out Cfg

	err := ecfg.ParsingFromCfgFile(path, key, &out)
	if err != nil {
		return out, fmt.Errorf("failed to parse %s from path %s: %s", key, path, err)
	}
	
	return out, nil
}