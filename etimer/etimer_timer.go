package etimer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/satori/go.uuid"
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
	jobs    map[string]*Job
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

func (t *Timer) AddJob(j *Job) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	j.js.mu.Lock()
	defer j.js.mu.Unlock()

	if j.timer != nil {
		return fmt.Errorf("job '%s' already exists in another timer", j.name)
	}

	o := t.jobs[j.name]
	if o != nil {
		return fmt.Errorf("job '%s' already exists", j.name)
	}

	interval := j.js.Interval()
	if j.js.pattern != "" {
		interval = time.Second
	}

	intervalTicksOfJob := int64(interval / t.options.Interval)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	nextTicks := atomic.LoadInt64(&t.ticks) + intervalTicksOfJob

	if j.js.pattern != "" {
		s, err := cron.ParseStandard(j.js.pattern)
		if err != nil {
			return err
		}
		j.sched = &scheduleCron{
			scheduleBasic: scheduleBasic{
				timer     : t,
				ticks     : intervalTicksOfJob,
				nextTicks_: nextTicks,
				js        : &j.js,
			},
			specSched: s,
		}
	} else {
		j.sched = &scheduleBasic{
			timer     : t,
			ticks     : intervalTicksOfJob,
			nextTicks_: nextTicks,
			js        : &j.js,
		}
	}

	j.timer = t
	t.queue.Push(j, nextTicks)

	return nil
}

func (t *Timer) AddInterval(ctx context.Context, interval time.Duration, cb JobFunc, singleton ...bool) *Job {
	j := createJob(JobOpts{
	  Name       : uuid.NewV1().String(),
		Ctx        : ctx,
		Interval   : interval,
		CB         : cb,
		IsSingleton: _getSingleton(singleton...),
		Times      : -1,
		status     : StatusReady,
	})
	t.AddJob(j)
	return j
}

func (t *Timer) AddCron(ctx context.Context, pattern string, cb JobFunc, singleton ...bool) (*Job, error) {
	j := createJob(JobOpts{
	  Name       : uuid.NewV1().String(),
		Ctx        : ctx,
		Pattern    : pattern,
		CB         : cb,
		IsSingleton: _getSingleton(singleton...),
		Times      : -1,
		status     : StatusReady,
	})
	err := t.AddJob(j)
	return j, err
}

func  (t *Timer) AllJobStates()[]*JobState{
	return nil
}

// return true in default
func _getSingleton(singleton ...bool) bool {
	if len(singleton) > 0 {
		return singleton[0]
	}
	return true
}


