package elog

import (
	"os"
	"testing"
)

func Test_SysLog(t *testing.T){

	syslog.Debugf("this is a elog sys dbg msg")
	syslog.Infof ("this is a elog sys inf msg")
	syslog.Warnf ("this is a elog sys wrn msg")
	syslog.Errorf("this is a elog sys err msg")

}

func TestLogLogic(t *testing.T) {
	os.Chdir("../")

	logger := Logger()

	logger.Log(ConsoleLevel(DebugLevel))
	logger.Log(ConsoleLevel(WarnLevel))
	logger.Log(ConsoleLevel(DebugLevel))

	if len(logger.consoleCores) != 2 {
		t.Errorf("expect 2, got: %d", len(logger.consoleCores))
	}
	if len(logger.consoleWriters) != 1 {
		t.Errorf("expect 1, got: %d", len(logger.consoleCores))
	}

	curFileCores   := len(logger.fileCores)
	curFileWriters := len(logger.fileWriters)

	logger.Log(FileLevel(FatalLevel))
	if len(logger.fileCores) != curFileCores + 1 {
		t.Errorf("expect %d, got %d", curFileCores + 1, len(logger.fileCores))
	}
	if len(logger.fileWriters) != curFileWriters {
		t.Errorf("expect %d, got %d", curFileWriters, len(logger.fileWriters))
	}
}

func TestLevelSetting(t *testing.T){

	os.Chdir("../")

	log := Log(NoFile(), NoConsole())
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should not output")
	log.Errorf("should not output")

	// test for console
	log = Log(NoFile(), ConsoleLevel(DebugLevel))
	log.Debugf("should output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(NoFile(), ConsoleLevel(InfoLevel))
	log.Debugf("should not output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(NoFile(), ConsoleLevel(WarnLevel))
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should output")
	log.Errorf("should output")

	// test for file write
	log = Log(Filename("level_test"), FileLevel(DebugLevel))
	log.Debugf("should output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(Filename("level_test"), FileLevel(InfoLevel))
	log.Debugf("should not output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(Filename("level_test"), FileLevel(WarnLevel))
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should output")
	log.Errorf("should output")
}

func TestDfLogger(t *testing.T){

	os.Chdir("../")

	log := Log()
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Filename("test"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Filename("test"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Filename("test")).Named("[test]")
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}

func TestYmlLogger(t *testing.T){

	os.Chdir("../")

	InitFromYml("./config/elog.yml")

	log := Log()
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log()
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Tag("[test1]"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Filename("test1"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Logger("log1").Log(Tag("[log1]"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Logger("log2").Log(Tag("[log2]"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Logger("log3").Log(Tag("[log3]"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Logger( "not_have", "log3").Log(Tag("[log3]"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}

