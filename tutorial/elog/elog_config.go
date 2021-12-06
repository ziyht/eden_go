package elog

import (
	"github.com/ziyht/eden_go/elog"
)

func runConfigTutorial() {

	elog.InitFromCfgFile("config/elog.yml")

	syslog.Infof("----------------- runConfigTutorial: default logger ---------------")
	log := elog.Logger().Log().Named("[default]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: nofile logger ---------------")
	log = elog.Logger("nofile").Log().Named("[nofile]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: setfile logger ---------------")
	log = elog.Logger("setfile").Log().Named("[setfile]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutorial: level_debug logger ---------------")
	log = elog.Logger("level_debug").Log().Named("[level_debug]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}