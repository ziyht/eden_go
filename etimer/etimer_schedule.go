package etimer

import (
	"sync/atomic"
	"time"
)

type schedule interface {
	nextTicks() int64
	calNextTicksAndStart(curTimerTicks int64, curTime time.Time)(reach bool, nextTicks int64, nextStart time.Time)
	commitNextTicks(nextTicks int64)
	commitNextStart(nextStart time.Time)

	doCheckTicksAndTime(curTimerTicks int64, curTime time.Time)bool
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

func (s *scheduleBasic) commitNextTicks(nextTicks int64) {
	atomic.StoreInt64(&s.nextTicks_, nextTicks)
}

func (s *scheduleBasic) commitNextStart(nextStart time.Time) {
	s.js.nextStart = nextStart
}

// only returns value when reach == true
func (s *scheduleBasic) calNextTicksAndStart(curTimerTicks int64, curTime time.Time)(reach bool, nextTicks int64, nextStart time.Time){
	// Ticks check.
	if curTimerTicks < atomic.LoadInt64(&s.nextTicks_) {
		return
	}

	return true, curTimerTicks + s.ticks, curTime.Add(s.timer.options.Interval * time.Duration(s.ticks))
}

// doCheckTicksAndTime checks the if job can run in given timer ticks or time,
// it returns true if the job need run else return false.
func (s *scheduleBasic) doCheckTicksAndTime(curTimerTicks int64, curTime time.Time) bool {
	reach, nt, ns := s.calNextTicksAndStart(curTimerTicks, curTime)

	if reach {
		s.commitNextTicks(nt)
		s.commitNextStart(ns)
	}

	return reach
}
