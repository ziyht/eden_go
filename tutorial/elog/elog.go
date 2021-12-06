package elog

import "eden_go/elog"

var syslog = elog.Logger().Log(elog.NoFile()).Named("[SYSLOG]")

func RunTutorail() {
  //runDefaultTutotil()
	runConfigTutotil()
}

