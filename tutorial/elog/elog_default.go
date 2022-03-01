package elog

import (
	"github.com/ziyht/eden_go/elog"
)

func runDefaultTutorial() {

	logger := elog.Logger()
	log := logger.Log()


	syslog.Infof("----------------- runDefaultTutotial 1 ---------------")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
	//log.Fatalf( "out put debug")

	syslog.Infof("----------------- runDefaultTutotial 2 ---------------")
	// using options to change log functions
	log = logger.Log(elog.Opt().ConsoleLevel(elog.LEVEL_DEBUG))
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}