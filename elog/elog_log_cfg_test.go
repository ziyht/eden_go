package elog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func parsingCfgFromStr(str string)(*Cfg, error) {

	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	return parsingCfg(&dfCfg, v, "elog", "")
}

func Test_parsingCfg(t *testing.T){


	var err error

	cfg, _ := parsingCfgFromStr("elog:\n dir: testval")
	if cfg.Dir != "testval" { t.Fatalf("val not match") }

	cfg, _  = parsingCfgFromStr("elog:\n group: testval")
	if cfg.Group != "testval" { t.Fatalf("val not match") }

	cfg, _  = parsingCfgFromStr("elog:\n filename: testval")
	if cfg.FileName != "testval" { t.Fatalf("val not match") }

	cfg, err  = parsingCfgFromStr("elog:\n max_size: testval")
	if cfg != nil { t.Fatalf("val not match") }
	if err == nil { t.Fatalf("err should occor") }
	cfg, err  = parsingCfgFromStr("elog:\n max_size: -1")
	if cfg != nil { t.Fatalf("val not match") }
	if err == nil { t.Fatalf("err should occor") }
	cfg, err  = parsingCfgFromStr("elog:\n max_size: 0")
	if cfg.MaxSize != 0 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }
	cfg, err  = parsingCfgFromStr("elog:\n max_size: 1024")
	if cfg.MaxSize != 1024 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }

	cfg, err  = parsingCfgFromStr("elog:\n max_backups: 0")
	if cfg.MaxBackups != 0 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }
	cfg, err  = parsingCfgFromStr("elog:\n max_backups: 1024")
	if cfg.MaxBackups  != 1024 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }

	cfg, err  = parsingCfgFromStr("elog:\n max_age: 0")
	if cfg.MaxAge != 0 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }
	cfg, err  = parsingCfgFromStr("elog:\n max_age: 1024")
	if cfg.MaxAge  != 1024 { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }

	cfg, err  = parsingCfgFromStr("elog:\n compress: true")
	if cfg.Compress != true { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }
	cfg, err  = parsingCfgFromStr("elog:\n compress: false")
	if cfg.Compress  != false { t.Fatalf("val not match") }
	if err != nil { t.Fatalf("err should not occor") }

	cfg, _  = parsingCfgFromStr("elog:\n color: false")
	if cfg.ConsoleColor  != false { t.Fatalf("val not match") }
	if cfg.FileColor     != false { t.Fatalf("val not match") }
	cfg, _  = parsingCfgFromStr("elog:\n color: true")
	if cfg.ConsoleColor  != true { t.Fatalf("val not match") }
	if cfg.FileColor     != true { t.Fatalf("val not match") }
	cfg, _  = parsingCfgFromStr("elog:\n color: [false, true]")
	if cfg.ConsoleColor  != false { t.Fatalf("val not match") }
	if cfg.FileColor     != true { t.Fatalf("val not match") }
	cfg, _  = parsingCfgFromStr("elog:\n color: [true, false]")
	if cfg.ConsoleColor  != true { t.Fatalf("val not match") }
	if cfg.FileColor     != false { t.Fatalf("val not match") }
}

func compareCfg(a *Cfg, b *Cfg) error {

	if a.Dir               !=  b.Dir              { return fmt.Errorf("Dir not match") }
	if a.Group             !=  b.Group            { return fmt.Errorf("Group not match") }
	if a.FileName          !=  b.FileName         { return fmt.Errorf("FileName not match") }
	if a.MaxSize           !=  b.MaxSize          { return fmt.Errorf("MaxSize not match") }
	if a.MaxBackups        !=  b.MaxBackups       { return fmt.Errorf("MaxBackups not match") }
	if a.MaxAge            !=  b.MaxAge           { return fmt.Errorf("MaxAge not match") }
	if a.Compress          !=  b.Compress         { return fmt.Errorf("Compress not match") }
	if a.ConsoleLevel      !=  b.ConsoleLevel     { return fmt.Errorf("ConsoleLevel not match") }
	if a.FileLevel         !=  b.FileLevel        { return fmt.Errorf("FileLevel not match") }
	if a.ConsoleStackLevel !=  b.ConsoleStackLevel{ return fmt.Errorf("ConsoleStackLevel not match") }
	if a.FileStackLevel    !=  b.FileStackLevel   { return fmt.Errorf("FileStackLevel not match") }
	if a.ConsoleColor      !=  b.ConsoleColor     { return fmt.Errorf("ConsoleColor not match") }	
	if a.FileColor         !=  b.FileColor        { return fmt.Errorf("FileColor not match") }	

	return nil
}

func Test_parsingDfCfg(t *testing.T){

	var testConfigContent = `
elog:
  dir           : logs                    # default logs
  group         : <HOSTNAME>              # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename      : <LOG_NAME>              # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File
  max_size      : 1                       # default 100, unit MB
  max_backups   : 1                       # default 7
  max_age       : 1                       # default 7
  compress      : true                    # default false
  level         : info                    # default info , debug, [0] for console, [1] for file, valid value is [debug, info, warn, error, fatal, panic]
  stack_level   : [warn, debug ]          # default fatal, warn , [0] for console, [1] for file, valid value is [debug, info, warn, error, fatal, panic]
  color         : false                   # default true , true , [0] for console, [1] for file

  none_reset: {}

  reset_goup:
    group   : group2

  reset_filename:
    filename: filename2

  reset_max_size:
    max_size: 11

  reset_max_backups:
    max_backups: 11

  reset_max_age:
    max_age: 11

  reset_compress:
    compress: false

  reset_level:
    level: [error, error]

  reset_stack_level:
    stack_level: [error, error]	

  reset_color:
    color: true
`

	cfgs, err := parsingCfgsFromStr(testConfigContent, "yml")
	if err != nil {  t.Fatalf("err occured: %s", err) }

	dfcfg := dfCfg

	if dfcfg.Dir               != "logs" { t.Fatalf("val not match") }
	if dfcfg.Group             != "<HOSTNAME>" { t.Fatalf("val not match") }
	if dfcfg.FileName          != "<LOG_NAME>" { t.Fatalf("val not match") }
	if dfcfg.MaxSize           != 1 { t.Fatalf("val not match") }
	if dfcfg.MaxBackups        != 1 { t.Fatalf("val not match") }
	if dfcfg.MaxAge            != 1 { t.Fatalf("val not match") }
	if dfcfg.Compress          != true { t.Fatalf("val not match") }
	if dfcfg.ConsoleLevel      != "info" { t.Fatalf("val not match") }
	if dfcfg.FileLevel         != "info" { t.Fatalf("val not match") }
	if dfcfg.ConsoleStackLevel != "warn" { t.Fatalf("val not match") }
	if dfcfg.FileStackLevel    != "debug" { t.Fatalf("val not match") }
	if dfcfg.ConsoleColor      != false { t.Fatalf("val not match") }

	cfg := cfgs.Cfgs["filename"]
	if cfg != nil { t.Fatalf("should be nil") }

	expect := dfcfg
  cfg = cfgs.Cfgs["none_reset"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }

	expect = dfcfg
	expect.Group = "group2"
  cfg = cfgs.Cfgs["reset_goup"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }

	expect = dfcfg
	expect.FileName = "filename2"
  cfg = cfgs.Cfgs["reset_filename"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }

	expect = dfcfg
	expect.MaxSize = 11
  cfg = cfgs.Cfgs["reset_max_size"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.MaxBackups = 11
  cfg = cfgs.Cfgs["reset_max_backups"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.MaxAge = 11
  cfg = cfgs.Cfgs["reset_max_age"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.Compress = false
  cfg = cfgs.Cfgs["reset_compress"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.ConsoleLevel = "error"
	expect.FileLevel    = "error"
  cfg = cfgs.Cfgs["reset_level"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.ConsoleStackLevel = "error"
	expect.FileStackLevel    = "error"
  cfg = cfgs.Cfgs["reset_stack_level"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	

	expect = dfcfg
	expect.ConsoleColor = true
	expect.FileColor    = true
  cfg = cfgs.Cfgs["reset_color"]
	if cfg == nil { t.Fatalf("cfg should not be nil") }
	err = compareCfg(cfg, &expect)
	if err != nil { t.Fatalf("%s", err) }	
}