package etimer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ziyht/eden_go/eerr"
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
			t.runJob(j)
		}
		
		// Status check: push back or ignore it.
		if j.js.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(j, j.sched.nextTicks())
		}
	}
}

func (t *Timer) submitJob(j *Job) {


}

func (t *Timer) runJob(j *Job) {
	go func() {
		start := time.Now()

		j.js.addRunning(start)

		defer func() {
			end   := time.Now()
			if exception := recover(); exception != nil {
				if exception != panicExit {
					if e, ok := exception.(error); ok && eerr.HasStack(e) {
					  j.js.addRunningOver(start, end, e)
						panic(e)
						
					} else {
						e := eerr.Newf(`exception recovered: %+v`, exception)
						j.js.addRunningOver(start, end, e)
						panic(e)
					}
				} else {
					j.Close()
					j.js.addRunningOver(start, end, nil)
					return
				}
			}
			if j.js.Status() == StatusRunning {
				j.js.setStatus(StatusReady)
			}
		}()

		j.js.setStart(start)
		err   := j.cb(j)
		end   := time.Now()
		j.js.addRunningOver(start, end, err)
	}()
}
