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

	logger.Log(Opt().ConsoleLevel(LEVEL_DEBUG))
	logger.Log(Opt().ConsoleLevel(LEVEL_WARN))
	logger.Log(Opt().ConsoleLevel(LEVEL_DEBUG))

	// if len(logger.consoleCores) != 2 {
	// 	t.Errorf("expect 2, got: %d", len(logger.consoleCores))
	// }
	// if len(logger.consoleWriters) != 1 {
	// 	t.Errorf("expect 1, got: %d", len(logger.consoleCores))
	// }

	// curFileCores   := len(logger.fileCores)
	// curFileWriters := len(logger.fileWriters)

	// logger.Log(FileLevel(LEVEL_FATAL))
	// if len(logger.fileCores) != curFileCores + 1 {
	// 	t.Errorf("expect %d, got %d", curFileCores + 1, len(logger.fileCores))
	// }
	// if len(logger.fileWriters) != curFileWriters {
	// 	t.Errorf("expect %d, got %d", curFileWriters, len(logger.fileWriters))
	// }
}

func TestLevelSetting(t *testing.T){

	os.Chdir("../")

	log := Log(Opt().NoFile().NoConsole())
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should not output")
	log.Errorf("should not output")

	// test for console
	log = Log(Opt().NoFile().ConsoleLevel(LEVEL_DEBUG))
	log.Debugf("should output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(Opt().NoFile().ConsoleLevel(LEVEL_INFO))
	log.Debugf("should not output")
	log.Infof ("should output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(Opt().NoFile().ConsoleLevel(LEVEL_WARN))
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should output")
	log.Errorf("should output")

	log = Log(Opt().NoFile().ConsoleLevel(LEVEL_ERROR))
	log.Debugf("should not output")
	log.Infof ("should not output")
	log.Warnf ("should not output")
	log.Errorf("should output")

	// test for file write
	log = Log(Opt().Filename("level_debug_test").FileLevel(LEVEL_DEBUG))
	log.Debugf("should write")
	log.Infof ("should write")
	log.Warnf ("should write")
	log.Errorf("should write")

	log = Log(Opt().Filename("level_info_test").FileLevel(LEVEL_INFO))
	log.Debugf("should not write")
	log.Infof ("should write")
	log.Warnf ("should write")
	log.Errorf("should write")

	log = Log(Opt().Filename("level_warn_test").FileLevel(LEVEL_WARN))
	log.Debugf("should not write")
	log.Infof ("should not write")
	log.Warnf ("should write")
	log.Errorf("should output")

	log = Log(Opt().Filename("level_error_test").FileLevel(LEVEL_ERROR))
	log.Debugf("should not write")
	log.Infof ("should not write")
	log.Warnf ("should not write")
	log.Errorf("should write")
}

func TestStackLevelSetting(t *testing.T){

	os.Chdir("../")

	// test for console
	log := Log(Opt().NoFile().ConsoleLevel(LEVEL_DEBUG).ConsoleStackLevel(LEVEL_INFO))
	log.Debugf("should output")
	log.Infof ("should output")
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

	log = Log(Opt().Filename("test"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Opt().Filename("test"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

}

func TestDfLoggerTags(t *testing.T){

	os.Chdir("../")

	logger := Logger()
	logger.option.Tags("first")

	log := Log(Opt().Filename("DfLoggerTags").Tags("tag1"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Opt().Filename("DfLoggerTags").Tags("tag1", "tag2"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Opt().Filename("DfLoggerTags").Tags("tag1", "", "tag3"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")

	log = Log(Opt().Filename("DfLoggerTags").ReTags("tag1").Tags("tag4"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}

func TestYmlLogger(t *testing.T){

	os.Chdir("../")

	InitFromFile("./config/elog.yml")

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

	log = Log(Opt().Filename("test1"))
	log.Debugf("dbg")
	log.Infof ("inf")
	log.Warnf ("wrn")
	log.Errorf("err")
}

