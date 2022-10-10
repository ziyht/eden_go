package etimer

import (
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

type scheduleCron struct {
  scheduleBasic
	specSched      cron.Schedule
}

func (s *scheduleCron) doCheckTicksAndTime(curTimerTicks int64, curTime time.Time) bool {
	needRun := false

	if curTime.After(s.js.nextStart){
		s.js.nextStart = s.specSched.Next(curTime)
		needRun = true
	}
	ticks := int64(s.js.nextStart.Sub(curTime) / s.timer.options.Interval)
	if ticks <= 0 {
		ticks = 1
	} else if ticks > s.ticks {
		ticks = s.ticks
	}
	atomic.StoreInt64(&s.nextTicks_, curTimerTicks + int64(ticks))

	if !needRun {
		return false
	}
	
	// Perform job checking.
	switch s.js.Status() {
	  case StatusRunning: if s.js.IsSingleton() { return false }
	  case StatusReady  : if !s.js.setStatusCas(StatusReady, StatusRunning) { return false }
	  case StatusStopped: return false
	  case StatusClosed : return false
	}	

	if s.js.times > 0 {
		leftRunTimes := atomic.AddInt64(&s.js.leftTimes, -1)
		if leftRunTimes < 0 {
			s.js.setStatus(StatusStopped)
			return false
		}
	}
	return true


}



