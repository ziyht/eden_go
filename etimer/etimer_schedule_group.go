package etimer

import (
	"time"
)


var nilTime time.Time

type scheduleGroup struct {
	queue   *priorityQueue
}

func (sg *scheduleGroup) nextTicks() int64{
	return sg.queue.nextPriority
}

func (sg *scheduleGroup) commitNextTicks(nextTicks int64) {
}

func (sg *scheduleGroup) commitNextStart(nextStart time.Time) {
}

func (sg *scheduleGroup) calNextTicksAndStart(curTimerTicks int64, curTime time.Time)(reach bool, nextTicks int64, nextStart time.Time){
	return
}

func (sg *scheduleGroup) doCheckTicksAndTime(curTimerTicks int64, curTime time.Time) bool {
	var (
		value     interface{}
		reach     bool
		nextStart time.Time
		s         schedule
	)

	for {
		value = sg.queue.Fetch()
		if value == nil {
			break
		}

		s = value.(schedule)
		reach, nt, ns := s.calNextTicksAndStart(curTimerTicks, curTime)
		if !reach {
			break
		}

		reach = true
		if nextStart == nilTime {
			nextStart = ns
		} else if nextStart.After(ns) {
			nextStart = ns
		}

		sg.queue.Pop()
		s.commitNextTicks(nt)
		sg.queue.Push(s, s.nextTicks())
	}

	if reach {
		s.commitNextStart(nextStart)
	}

	return reach
}


