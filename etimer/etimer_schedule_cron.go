package etimer

import (
	"time"

	"github.com/robfig/cron/v3"
)

type scheduleCron struct {
  scheduleBasic
	specSched      cron.Schedule
}

// only returns value when reach == true
func (s *scheduleCron) calNextTicksAndStart(curTimerTicks int64, curTime time.Time)(reach bool, nextTicks int64, nextStart time.Time){
	// time check.
	if curTime.Before(s.js.nextStart){
		return
	}

	nextStart = s.specSched.Next(curTime)
	ticks := int64(nextStart.Sub(curTime) / s.timer.options.Interval)
	if ticks <= 0 {
		ticks = 1
	} else if ticks > s.ticks {
		ticks = s.ticks
	}

	return true, curTimerTicks + ticks, nextStart
}

func (s *scheduleCron) doCheckTicksAndTime(curTimerTicks int64, curTime time.Time) bool {
	reach, nt, ns := s.calNextTicksAndStart(curTimerTicks, curTime)

	if reach {
		s.commitNextTicks(nt)
		s.commitNextStart(ns)
	}

	return reach
}



