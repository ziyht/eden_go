package elog

import "github.com/ziyht/eden_go/elog"

var syslog = elog.Logger().Log(elog.Opt().NoFile()).Named("[SYSLOG]")

func RunTutorail() {
  runDefaultTutorial()
	runConfigTutorial()
}

