package etimer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// internalPanic is the custom panic for internal usage.
type internalPanic string


// TimerOptions is the configuration object for Timer.
type TimerOptions struct {
	Interval time.Duration // Interval is the interval escaped of the timer.
}

// Timer is the timer manager, which uses ticks to calculate the timing interval.
type Timer struct {
	mu      sync.RWMutex
	queue   *priorityQueue // queue is a priority queue based on heap structure.
	status  int32          // status is the current timer status.
	ticks   int64          // ticks is the proceeded interval number by the timer.
	options TimerOptions   // timer options is used for timer configuration.
}

func newTimer(options ...TimerOptions) *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: StatusRunning,
		ticks:  0,
	}
	if len(options) > 0 {
		t.options = options[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// AddJob adds a timing job to the timer with detailed parameters.
//
// The parameter `interval` specifies the running interval of the job.
//
// The parameter `singleton` specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when it's a singleton mode job.
//
// The parameter `times` specifies limit for the job running times, which means the job
// exits if its run times exceeds the `times`.
//
// The parameter `status` specifies the job status when it's firstly added to the timer.
func (t *Timer) AddJob(ctx context.Context, interval time.Duration, cb JobFunc, isSingleton bool, times int64, status int32) *Job {
	return t.createJob(createJobInput{
		Ctx:         ctx,
		Interval:    interval,
		CB :         cb,
		IsSingleton: isSingleton,
		Times:       times,
		Status:      status,
	})
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(ctx context.Context, interval time.Duration, cb JobFunc) *Job {
	return t.createJob(createJobInput{
		Ctx:         ctx,
		Interval:    interval,
		CB :         cb,
		IsSingleton: true,
		Times:       -1,
		Status:      StatusReady,
	})
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
			case <-timerIntervalTicker.C:
				// Check the timer status.
				switch atomic.LoadInt32(&t.status) {
				case StatusRunning:
					// Timer proceeding.
					if currentTimerTicks = atomic.AddInt64(&t.ticks, 1); currentTimerTicks >= t.queue.NextPriority() {
						t.proceed(currentTimerTicks)
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
func (t *Timer) proceed(currentTimerTicks int64) {
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
		if jobNextTicks := atomic.LoadInt64(&j.nextTicks) ; currentTimerTicks < jobNextTicks {
			break
		}
		t.queue.Pop()
		// It checks the job running requirements and then does asynchronous running.
		j.doCheckAndRunByTicks(currentTimerTicks)
		// Status check: push back or ignore it.
		if j.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(j, atomic.LoadInt64(&j.nextTicks))
		}
	}
}

type createJobInput struct {
	Ctx         context.Context
	Interval    time.Duration
	CB          JobFunc
	IsSingleton bool
	Times       int64
	Status      int32
}

// createJob creates and adds a timing job to the timer.
func (t *Timer) createJob(in createJobInput) *Job {
	var (
		intervalTicksOfJob = int64(in.Interval / t.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	var (
		nextTicks = atomic.LoadInt64(&t.ticks) + intervalTicksOfJob
		j     = &Job{
			cb:          in.CB,
			ctx:         in.Ctx,
			timer:       t,
			ticks:       intervalTicksOfJob,
			nextTicks:   nextTicks,
		}
	)
	j.state.setStatus(in.Status)
	j.SetSingleton(in.IsSingleton)
	if in.Times > 0 {
		j.SetTimes(in.Times)
	}

	t.queue.Push(j, nextTicks)
	return j
}
