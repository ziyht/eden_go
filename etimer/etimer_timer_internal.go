package etimer

import (
	"fmt"
	"sync/atomic"
	"time"
)

func statusString(s int32) string {

	switch s {
		case StatusReady  : return "StatusReady"
		case StatusRunning: return "StatusRunning"
		case StatusStopped: return "StatusStopped"
		case StatusClosed : return "StatusClosed"
	}

	return fmt.Sprintf("StatusUnknown(%d)", s)
}

// loop starts the ticker using a standalone goroutine.
func (t *Timer) loop() {
	go func() {
		var (
			currentTimerTicks   int64
			timerIntervalTicker = time.NewTicker(t.options.Interval)
		)
		defer timerIntervalTicker.Stop()
		for {
			select {
			case now := <-timerIntervalTicker.C:
				// Check the timer status.
				switch atomic.LoadInt32(&t.status) {
				case StatusRunning:
					// Timer proceeding.
					if currentTimerTicks = atomic.AddInt64(&t.ticks, 1); currentTimerTicks >= t.queue.NextPriority() {
						t.proceed(currentTimerTicks, now)
					}

				case StatusStopped:
					// Do nothing.

				case StatusClosed:
					// Timer exits.
					return
				}
			}
		}
	}()
}

// proceed function proceeds the timer job checking and running logic.
func (t *Timer) proceed(curTimerTicks int64, curTime time.Time) {
	var (
		value interface{}
	)
	for {
		value = t.queue.Fetch()
		if value == nil {
			break
		}
		j := value.(*Job)
		// It checks if it meets the ticks' requirement.
		if curTimerTicks < j.sched.nextTicks() {
			break
		}
		t.queue.Pop()
		// It checks the job running requirements and then does asynchronous running.
		if j.sched.doCheckTicksAndTime(curTimerTicks, curTime){
			t.submitJob(j)
		}
		
		// Status check: push back or ignore it.
		if j.js.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(j, j.sched.nextTicks())
		}
	}
}

func (t *Timer) submitJob(j *Job) {
	dfRunner.run(j)
}

func (t *Timer) runJob(j *Job) {
	go j.exec_once()
}
