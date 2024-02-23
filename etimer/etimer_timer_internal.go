package etimer

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/ziyht/eden_go/eerr"
)

func statusString(s int32) string {

	switch s {
		case StatusWaiting: return "Waiting"
		case StatusRunning: return "Running"
		case StatusPending: return "Pending"
		case StatusStopped: return "Stopped"
		case StatusClosed : return "Closed"
	}

	return fmt.Sprintf("StatusUnknown(%d)", s)
}

func (t *Timer) parsingScheduleForJob(j *Job) error {
	intervalTicksOfJob := int64(j.js.Interval() / t.options.Interval)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	nextTicks := atomic.LoadInt64(&t.ticks) + intervalTicksOfJob

	if j.js.pattern == "" {
		j.sched = &scheduleBasic{
			timer     : t,
			ticks     : intervalTicksOfJob,
			nextTicks_: nextTicks,
			js        : &j.js,
		}
		return nil
	}

	var validPs []string 
	ps := strings.Split(j.js.pattern, ",")
	for _, p := range ps {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		validPs = append(validPs, p)
	}

	if len(validPs) == 0 {
		return eerr.Newf("invalid pattern(%s), can not find any valid pattern", j.js.pattern)
	}

	var sched schedule
	var sg *scheduleGroup

	for _, vp := range validPs {
		s, err := cron.ParseStandard(vp)
		if err != nil {
			return err
		}

		sched = &scheduleCron{
			scheduleBasic: scheduleBasic{
				timer     : t,
				ticks     : intervalTicksOfJob,
				nextTicks_: nextTicks,
				js        : &j.js,
			},
			specSched: s,
		}

		if len(validPs) == 1 {
			break
		} else if sg == nil {
			sg = &scheduleGroup{
				queue: newPriorityQueue(),
			}
		}

		sg.queue.Push(sched, sched.nextTicks())
		sched = sg
	}

	j.sched = sched 

	return nil
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

			// Perform job checking.
			switch j.js.Status() {
				case StatusStopped: return 
				case StatusClosed : return 
				default:
					fmt.Printf("submit '%s'\n", j.name)
					t.submitJob(j)
			}
		}
		
		// Status check: push back or ignore it.
		if j.js.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(j, j.sched.nextTicks())
		}
	}
}

func (t *Timer) submitJob(j *Job) {
	dfRunner.submit(j)
}
