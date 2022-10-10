package etimer

import (
	"sync/atomic"
	"time"
)

type schedule interface {
	doCheckTicksAndTime(curTimerTicks int64, curTime time.Time)bool
	nextTicks() int64
}

type scheduleBasic struct {
	timer       *Timer
	ticks       int64           // The job runs every tick.
	nextTicks_  int64           // Next run ticks of the job.
	
	js          *JobState
}

func (s *scheduleBasic) nextTicks() int64 {
	return atomic.LoadInt64(&s.nextTicks_)
}

func (s *scheduleBasic) CheckLimits() bool {
	if s.js.times <= 0 {
		return true
	}

	leftRunTimes := atomic.AddInt64(&s.js.leftTimes, -1)
	return leftRunTimes >= 0
}

// doCheckTicksAndTime checks the if job can run in given timer ticks or time,
// it returns true if the job need run else return false.
func (s *scheduleBasic) doCheckTicksAndTime(curTimerTicks int64, curTime time.Time) bool {

	// Ticks check.
	if curTimerTicks < atomic.LoadInt64(&s.nextTicks_) {
		return false
	}
	atomic.StoreInt64(&s.nextTicks_, curTimerTicks + s.ticks)
	s.js.nextStart = curTime.Add(s.timer.options.Interval * time.Duration(s.ticks))
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
