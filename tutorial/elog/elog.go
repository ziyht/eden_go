package elog

import "eden/elog"

var syslog = elog.Logger().Log(elog.NoFile()).Named("[SYSLOG]")

func RunTutorail() {
  //runDefaultTutotil()
	runConfigTutotil()
}

