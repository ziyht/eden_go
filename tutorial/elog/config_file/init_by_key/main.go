package main

import "github.com/ziyht/eden_go/elog"

func main(){
	logger := elog.LoggerFromFile("../elog.yml", "elog.log1")
	log1 := logger.Log()
  log1.Debugf("output debug")  // in default setting, it will not output
	log1.Infof( "output info")
	log1.Warnf( "output warn")
	log1.Errorf("output error")
}

