package main

import (
	"github.com/ziyht/eden_go/elog"
)

func main() {

	logger := elog.Logger()
	syslog := logger.Log(elog.Opt().NoFile()).Named("[SYSLOG]")
	log := logger.Log()

	syslog.Infof("----------------- runDefaultTutotial 1 ---------------")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
	//log.Fatalf( "out put debug")

	// using options to change log functions
	log = logger.Log(elog.Opt().ConsoleLevel(elog.LEVEL_DEBUG))
	syslog.Infof("----------------- runDefaultTutotial 2 ---------------")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}

