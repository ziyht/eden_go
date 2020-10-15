package elog

import (
	"testing"
)

func Test_SysLog(t *testing.T){

	syslog.Debugf("this is a elog sys dbg msg")
	syslog.Infof ("this is a elog sys inf msg")
	syslog.Warnf ("this is a elog sys wrn msg")
	syslog.Errorf("this is a elog sys err msg")

}

func TestDfLogger(t *testing.T){

	log := DfLog()
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLog("test")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLog("test")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLog("test").Named("[test]")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}

func TestYmlLogger(t *testing.T){

	// path, _ := filepath.Abs("../config/elog.yml")

	InitFromYml("../config/elog.yml")

	log := DfLog()
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLog("")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLog("test1")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLogger("log1").GetLog().Named("[log1]")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = GetLogger("log2").GetLog().Named("[log2]")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}
