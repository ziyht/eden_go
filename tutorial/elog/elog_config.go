package elog

import (
	"github.com/ziyht/eden_go/elog"
)

func runConfigTutotil() {

	elog.InitFromFile("config/elog.yml")

	syslog.Infof("----------------- runConfigTutotil: default logger ---------------")
	log := elog.Logger().Log().Named("[default]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutotil: nofile logger ---------------")
	log = elog.Logger("nofile").Log().Named("[nofile]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutotil: setfile logger ---------------")
	log = elog.Logger("setfile").Log().Named("[setfile]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	syslog.Infof("----------------- runConfigTutotil: level_debug logger ---------------")
	log = elog.Logger("level_debug").Log().Named("[level_debug]")
	log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}