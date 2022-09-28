package main

import "github.com/ziyht/eden_go/elog"


var	syslog = elog.Logger().Log(elog.Opt().NoFile()).Named("[SYSLOG]")

func main() {

	elog.InitFromFile("../elog.yml")

	syslog.Infof("----------------- runConfigTutorial: default logger ---------------")
	log := elog.Logger().Log().Named("[default]")
	log.Debugf("output debug")
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: log1 logger ---------------")
	log = elog.Logger("log1").Log()
	log.Debugf("output debug") 
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: log2 logger ---------------")
	log = elog.Logger("log2").Log()
	log.Debugf("output debug") 
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: multi_file logger ---------------")
	log = elog.Logger("multi_file").Log()
	log.Debugf("output debug") 
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: only_console logger ---------------")
	log = elog.Logger("only_console").Log()
	log.Debugf("output debug") 
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}


