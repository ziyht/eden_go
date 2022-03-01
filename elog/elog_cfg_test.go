package elog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func parsingCfgFromStr(str string)(*LoggerCfg, error) {

	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	return parsingLoggerCfg(&dfCfg, v, "elog", "", nil)
}

func Test_parsingCfg(t *testing.T){


	//var err error

	// cfg, _ := parsingCfgFromStr("elog:\n dir: testval")
	// if cfg.Dir != "testval" { t.Fatalf("val not match") }

	// cfg, _  = parsingCfgFromStr("elog:\n group: testval")
	// if cfg.Group != "testval" { t.Fatalf("val not match") }

	// cfg, _  = parsingCfgFromStr("elog:\n filename: testval")
	// if cfg.FileName != "testval" { t.Fatalf("val not match") }

	// cfg, err  = parsingCfgFromStr("elog:\n max_size: testval")
	// if cfg != nil { t.Fatalf("val not match") }
	// if err == nil { t.Fatalf("err should occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n max_size: -1")
	// if cfg != nil { t.Fatalf("val not match") }
	// if err == nil { t.Fatalf("err should occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n max_size: 0")
	// if cfg.MaxSize != 0 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n max_size: 1024")
	// if cfg.MaxSize != 1024 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }

	// cfg, err  = parsingCfgFromStr("elog:\n max_backups: 0")
	// if cfg.MaxBackup != 0 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n max_backups: 1024")
	// if cfg.MaxBackup  != 1024 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }

	// cfg, err  = parsingCfgFromStr("elog:\n max_age: 0")
	// if cfg.MaxAge != 0 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n max_age: 1024")
	// if cfg.MaxAge  != 1024 { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }

	// cfg, err  = parsingCfgFromStr("elog:\n compress: true")
	// if cfg.Compress != true { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }
	// cfg, err  = parsingCfgFromStr("elog:\n compress: false")
	// if cfg.Compress  != false { t.Fatalf("val not match") }
	// if err != nil { t.Fatalf("err should not occor") }

	// cfg, _  = parsingCfgFromStr("elog:\n color: false")
	// if cfg.ConsoleColor  != ColorOff { t.Fatalf("val not match") }
	// if cfg.FileColor     != ColorOff { t.Fatalf("val not match") }
	// cfg, _  = parsingCfgFromStr("elog:\n color: true")
	// if cfg.ConsoleColor  != ColorOn { t.Fatalf("val not match") }
	// if cfg.FileColor     != ColorOn { t.Fatalf("val not match") }
	// cfg, _  = parsingCfgFromStr("elog:\n color: [false, true]")
	// if cfg.ConsoleColor  != ColorOff { t.Fatalf("val not match") }
	// if cfg.FileColor     != ColorOn { t.Fatalf("val not match") }
	// cfg, _  = parsingCfgFromStr("elog:\n color: [true, false]")
	// if cfg.ConsoleColor  != ColorOn { t.Fatalf("val not match") }
	// if cfg.FileColor     != ColorOff { t.Fatalf("val not match") }
}

func compareCfg(a *Cfg, b *Cfg) error {

	if a.Dir               !=  b.Dir              { return fmt.Errorf("Dir not match") }
	if a.Group             !=  b.Group            { return fmt.Errorf("Group not match") }
	if a.FileName          !=  b.FileName         { return fmt.Errorf("FileName not match") }
	if a.MaxSize           !=  b.MaxSize          { return fmt.Errorf("MaxSize not match") }
	if a.MaxBackup         !=  b.MaxBackup        { return fmt.Errorf("MaxBackups not match") }
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

// 	var testConfigContent = `
// elog:

//   dir           : logs                # default logs
//   group         : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $Dir/$Group
//   filename      : <APP>_<LOG>         # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File
//   max_size      : 100                 # default 100, unit MB
//   max_backups   : 7                   # default 7
//   max_age       : 7                   # default 7
//   compress      : false               # default false
//   f_level       : debug               # default debug, for file, valid value is [debug, info, warn, error, fatal, panic]
//   f_stack_level : warn                # default warn , for file, valid value is [debug, info, warn, error, fatal, panic]
//   f_color       : false               # default false, for file
//   c_level       : info                # default info , for console, valid value is [debug, info, warn, error, fatal, panic]
//   c_stack_level : error               # default error, for console, valid value is [debug, info, warn, error, fatal, panic]
//   c_color       : true                # default true , for console

//   none_reset: []

//   reset_goup:
// 	- inherit: default.console
// 	- inherit: default.file
//     group  : group2

//   reset_filename:
//     filename: filename2

//   reset_max_size:
//     max_size: 11

//   reset_max_backups:
//     max_backups: 11

//   reset_max_age:
//     max_age: 11

//   reset_compress:
//     compress: false

//   reset_level:
//     level: [error, error]

//   reset_stack_level:
//     stack_level: [error, error]	

//   reset_color:
//     color: true
// `

	// cfgs, err := parsingCfgsFromStr(testConfigContent, "yml")
	// if err != nil {  t.Fatalf("err occured: %s", err) }

	// dfcfg := dfCfg

	// if dfcfg.Dir               != "logs" { t.Fatalf("val not match") }
	// if dfcfg.Group             != "<HOSTNAME>" { t.Fatalf("val not match") }
	// if dfcfg.FileName          != "<LOG_NAME>" { t.Fatalf("val not match") }
	// if dfcfg.MaxSize           != 1 { t.Fatalf("val not match") }
	// if dfcfg.MaxBackup         != 1 { t.Fatalf("val not match") }
	// if dfcfg.MaxAge            != 1 { t.Fatalf("val not match") }
	// if dfcfg.Compress          != true { t.Fatalf("val not match") }
	// if dfcfg.ConsoleLevel      != LEVEL_INFO { t.Fatalf("val not match") }
	// if dfcfg.FileLevel         != LEVEL_INFO { t.Fatalf("val not match") }
	// if dfcfg.ConsoleStackLevel != LEVEL_WARN { t.Fatalf("val not match") }
	// if dfcfg.FileStackLevel    != LEVEL_DEBUG { t.Fatalf("val not match") }
	// if dfcfg.ConsoleColor      != ColorOff { t.Fatalf("val not match") }

	// cfg := cfgs["filename"]
	// if cfg != nil { t.Fatalf("should be nil") }

	// expect := dfcfg
  // cfg = cfgs["none_reset"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }

	// expect = dfcfg
	// expect.Group = "group2"
  // cfg = cfgs["reset_goup"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }

	// expect = dfcfg
	// expect.FileName = "filename2"
  // cfg = cfgs["reset_filename"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }

	// expect = dfcfg
	// expect.MaxSize = 11
  // cfg = cfgs["reset_max_size"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.MaxBackups = 11
  // cfg = cfgs["reset_max_backups"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.MaxAge = 11
  // cfg = cfgs["reset_max_age"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.Compress = false
  // cfg = cfgs["reset_compress"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.ConsoleLevel = "error"
	// expect.FileLevel    = "error"
  // cfg = cfgs["reset_level"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.ConsoleStackLevel = "error"
	// expect.FileStackLevel    = "error"
  // cfg = cfgs["reset_stack_level"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	

	// expect = dfcfg
	// expect.ConsoleColor = true
	// expect.FileColor    = true
  // cfg = cfgs["reset_color"]
	// if cfg == nil { t.Fatalf("cfg should not be nil") }
	// err = compareCfg(cfg, &expect)
	// if err != nil { t.Fatalf("%s", err) }	
}