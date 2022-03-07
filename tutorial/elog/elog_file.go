package elog

import "github.com/ziyht/eden_go/elog"


func runFileTutorial(){

	logger := elog.LoggerFromFile("config/elog.yml", "elog.console")
	log1 := logger.Log()
  log1.Debugf("output debug")  // in default setting, it will not output
	log1.Infof( "output info")
	log1.Warnf( "output warn")
	log1.Errorf("output error")

}

