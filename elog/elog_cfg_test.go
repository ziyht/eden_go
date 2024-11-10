package elog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func parsingCfgFromStr(str string)(*LoggerCfg, error) {

	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	return parsingLoggerCfg(&dfLoggerCfg, v, "elog", "", nil)
}

func Test_parsingCfg(t *testing.T){

	var cfg *LoggerCfg
	var err error

	cfg, err = parsingCfgFromStr("elog:\n  dir: test_val")
	assert.Nil(t, err)
	assert.Equal(t, "file", cfg.cfgs[0].Name)
	assert.Equal(t, "test_val", cfg.cfgs[0].Dir)

	cfg, err = parsingCfgFromStr("elog:\n  group: test_val")
	assert.Nil(t, err)
	assert.Equal(t, "test_val", cfg.cfgs[0].Group)
	
	cfg, err  = parsingCfgFromStr("elog:\n  filename: test_val")
	assert.Nil(t, err)
	assert.Equal(t, "test_val", cfg.cfgs[0].FileName)

	_, err  = parsingCfgFromStr("elog:\n  max_size: test_val")
	assert.NotNil(t, err)

	_, err  = parsingCfgFromStr("elog:\n  max_size: -1")
	assert.NotNil(t, err)

	cfg, err  = parsingCfgFromStr("elog:\n  max_size: 0")
	assert.Nil(t, err)
	assert.Equal(t, 0, cfg.cfgs[0].MaxSize)

	cfg, err  = parsingCfgFromStr("elog:\n  max_size: 1024")
	assert.Nil(t, err)
	assert.Equal(t, 1024, cfg.cfgs[0].MaxSize)

	_, err  = parsingCfgFromStr("elog:\n  max_backup: -1")
	assert.NotNil(t, err)

	cfg, err  = parsingCfgFromStr("elog:\n  max_backup: 0")
	assert.Nil(t, err)
	assert.Equal(t, 0, cfg.cfgs[0].MaxBackup)

	cfg, err  = parsingCfgFromStr("elog:\n  max_backup: 1024")
	assert.Nil(t, err)
	assert.Equal(t, 1024, cfg.cfgs[0].MaxBackup)

	_, err  = parsingCfgFromStr("elog:\n  max_age: -1")
	assert.NotNil(t, err)
	cfg, err  = parsingCfgFromStr("elog:\n  max_age: 0")
	assert.Nil(t, err)
	assert.Equal(t, 0, cfg.cfgs[0].MaxAge)
	cfg, err  = parsingCfgFromStr("elog:\n  max_age: 1024")
	assert.Nil(t, err)
	assert.Equal(t, 1024, cfg.cfgs[0].MaxAge)

	_, err  = parsingCfgFromStr("elog:\n  compress: bad")
	assert.NotNil(t, err)
	cfg, err  = parsingCfgFromStr("elog:\n  compress: true")
	assert.Nil(t, err)
	assert.Equal(t, true, cfg.cfgs[0].Compress)
	cfg, err  = parsingCfgFromStr("elog:\n  compress: false")
	assert.Nil(t, err)
	assert.Equal(t, false, cfg.cfgs[0].Compress)
	
	_, err = parsingCfgFromStr("elog:\n  color: bad")
	assert.NotNil(t, err)
	cfg, err = parsingCfgFromStr("elog:\n  color: true")
	assert.Nil(t, err)
	assert.Equal(t, ColorOn, cfg.cfgs[0].Color)
	assert.Equal(t, ColorOn, cfg.cfgs[1].Color)

	cfg, err = parsingCfgFromStr("elog:\n  color: false")
	assert.Nil(t, err)
	assert.Equal(t, ColorOff, cfg.cfgs[0].Color)
	assert.Equal(t, ColorOff, cfg.cfgs[1].Color)

	cfg, err = parsingCfgFromStr("elog:\n  color: auto")
	assert.Nil(t, err)
	assert.Equal(t, ColorAuto, cfg.cfgs[0].Color)
	assert.Equal(t, ColorAuto, cfg.cfgs[1].Color)
}

func compareCfg(a *LoggerCfg, b *LoggerCfg) error {

	if len(a.cfgs) != len(b.cfgs) { return fmt.Errorf("cfgs length not match") }

	for i := 0; i < len(a.cfgs); i++ {
		if a.cfgs[i].file              !=  b.cfgs[i].file             { return fmt.Errorf("df[%d].Name not match", i) }
		if a.cfgs[i].Tag               !=  b.cfgs[i].Tag              { return fmt.Errorf("df[%d].Name not match", i) }
		if a.cfgs[i].Name              !=  b.cfgs[i].Name             { return fmt.Errorf("df[%d].Name not match", i) }
		if a.cfgs[i].Dir               !=  b.cfgs[i].Dir              { return fmt.Errorf("df[%d].Dir not match", i) }
		if a.cfgs[i].Group             !=  b.cfgs[i].Group            { return fmt.Errorf("df[%d].Group not match", i) }
		if a.cfgs[i].FileName          !=  b.cfgs[i].FileName         { return fmt.Errorf("df[%d].FileName not match", i) }
		if a.cfgs[i].MaxSize           !=  b.cfgs[i].MaxSize          { return fmt.Errorf("df[%d].MaxSize not match", i) }
		if a.cfgs[i].MaxBackup         !=  b.cfgs[i].MaxBackup        { return fmt.Errorf("df[%d].MaxBackups not match", i) }
		if a.cfgs[i].MaxAge            !=  b.cfgs[i].MaxAge           { return fmt.Errorf("df[%d].MaxAge not match", i) }
		if a.cfgs[i].Compress          !=  b.cfgs[i].Compress         { return fmt.Errorf("df[%d].Compress not match", i) }
		if a.cfgs[i].Level      			 !=  b.cfgs[i].Level     				{ return fmt.Errorf("df[%d].ConsoleLevel not match", i) }
		if a.cfgs[i].StackLevel        !=  b.cfgs[i].StackLevel       { return fmt.Errorf("df[%d].FileLevel not match", i) }
		if a.cfgs[i].Color             !=  b.cfgs[i].Color            { return fmt.Errorf("df[%d].Color not match", i) }
	}

	return nil
}

func Test_parsingDfCfg(t *testing.T){

	var testConfigContent = `
elog:

  dir           : logs                # default logs
  group         : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename      : <APP>_<LOG>         # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File
  max_size      : 100                 # default 100, unit MB
  max_backup    : 7                   # default 7
  max_age       : 7                   # default 7
  compress      : false               # default false
  f_level       : debug               # default debug, for file, valid value is [debug, info, warn, error, fatal, panic]
  f_slevel 			: warn                # default warn , for file, valid value is [debug, info, warn, error, fatal, panic]
  f_color       : auto                # default false, for file
  c_level       : info                # default info , for console, valid value is [debug, info, warn, error, fatal, panic]
  c_slevel 			: error               # default error, for console, valid value is [debug, info, warn, error, fatal, panic]
  c_color       : auto                # default true , for console

  none_reset: []

  reset_group:
  - inherit: default.file
    group  : group2
  - inherit: default.console

  reset_filename:
    filename: filename2

  reset_max_size:
    max_size: 11

  reset_max_backups:
    max_backup: 11

  reset_max_age:
    max_age: 11

  reset_compress:
    compress: false

  reset_level:
    c_level: error
    f_level: error

  reset_stack_level:
    c_slevel: error
    f_slevel: error

  reset_color:
    c_color: true
    f_color: true
`

	cfgs, err := parsingCfgsFromStr(testConfigContent, "yml", "elog")
	assert.Nil(t, err)

	dfcfg := dfLoggerCfg.genLoggerCfg()
	assert.Equal(t, true          , dfcfg.cfgs[0].file)
	assert.Equal(t, "file"				, dfcfg.cfgs[0].Name)
	assert.Equal(t, ""						, dfcfg.cfgs[0].Tag)
	assert.Equal(t, "logs"				, dfcfg.cfgs[0].Dir)
	assert.Equal(t, "<HOSTNAME>"	, dfcfg.cfgs[0].Group)
	assert.Equal(t, "<APP>_<LOG>"	, dfcfg.cfgs[0].FileName)
	assert.Equal(t, 100						, dfcfg.cfgs[0].MaxSize)
	assert.Equal(t, 7						  , dfcfg.cfgs[0].MaxBackup)
	assert.Equal(t, 7						  , dfcfg.cfgs[0].MaxAge)
	assert.Equal(t, false				  , dfcfg.cfgs[0].Compress)
	assert.Equal(t, LEVEL_DEBUG   , dfcfg.cfgs[0].Level)
	assert.Equal(t, LEVEL_WARN    , dfcfg.cfgs[0].StackLevel)
	assert.Equal(t, ColorAuto     , dfcfg.cfgs[0].Color)

	assert.Equal(t, false         , dfcfg.cfgs[1].file)
	assert.Equal(t, "console"		  , dfcfg.cfgs[1].Name)
	assert.Equal(t, ""						, dfcfg.cfgs[1].Tag)
	assert.Equal(t, ""    				, dfcfg.cfgs[1].Dir)
	assert.Equal(t, ""						, dfcfg.cfgs[1].Group)
	assert.Equal(t, ""						, dfcfg.cfgs[1].FileName)
	assert.Equal(t, 0							, dfcfg.cfgs[1].MaxSize)
	assert.Equal(t, 0							, dfcfg.cfgs[1].MaxBackup)
	assert.Equal(t, 0							, dfcfg.cfgs[1].MaxAge)
	assert.Equal(t, false				  , dfcfg.cfgs[1].Compress)
	assert.Equal(t, LEVEL_INFO    , dfcfg.cfgs[1].Level)
	assert.Equal(t, LEVEL_ERROR   , dfcfg.cfgs[1].StackLevel)
	assert.Equal(t, ColorAuto     , dfcfg.cfgs[1].Color)


	cfg := cfgs["filename"]
	assert.Nil(t, cfg)

	expect := dfcfg.Clone()
  cfg = cfgs["none_reset"]

	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].Group = "group2"
  cfg = cfgs["reset_group"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].FileName = "filename2"
  cfg = cfgs["reset_filename"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].MaxSize = 11
  cfg = cfgs["reset_max_size"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].MaxBackup = 11
  cfg = cfgs["reset_max_backups"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].MaxAge = 11
  cfg = cfgs["reset_max_age"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].Compress = false
  cfg = cfgs["reset_compress"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].Level = LEVEL_ERROR
	expect.cfgs[1].Level = LEVEL_ERROR
  cfg = cfgs["reset_level"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].StackLevel = LEVEL_ERROR
	expect.cfgs[1].StackLevel = LEVEL_ERROR
  cfg = cfgs["reset_stack_level"]
	assert.Nil(t, compareCfg(cfg, expect))

	expect = dfcfg.Clone()
	expect.cfgs[0].Color = ColorOn
	expect.cfgs[1].Color = ColorOn
  cfg = cfgs["reset_color"]
	assert.Nil(t, compareCfg(cfg, expect))
}