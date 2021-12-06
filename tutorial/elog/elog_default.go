package elog

import (
	"eden_go/elog"
	"fmt"
)

func runDefaultTutotil() {

	logger := elog.Logger()
	log := logger.Log()
	fmt.Printf("%v\n", logger.Cfg())

	syslog.Infof("----------------- runDefaultTutotil 1 ---------------")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
	//log.Fatalf( "out put debug")

	syslog.Infof("----------------- runDefaultTutotil 2 ---------------")
	// using options to change log functions
	log = logger.Log(elog.ConsoleLevel(elog.LEVEL_DEBUG))
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}