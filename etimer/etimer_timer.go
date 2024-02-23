package etimer

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	groups  map[string]*runner
}

func newTimer(options ...TimerOptions) *Timer {
	t := &Timer{
	  jobs  : map[string]*Job{},
		queue :  newPriorityQueue(),
		status: StatusRunning,
		ticks :  0,
		groups:  map[string]*runner{},
	}
	if len(options) > 0 {
		t.options = options[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// SetGroup - set a running group in timer
// 1. a group is a handle to limit the max runings jobs in the same time, 
//    each job can be set to a seperate group
// 2. a timer can have several groups
// 3. a group will be created with the name if it not exist, or it will be update by @param max for 
func (t *Timer) SetGroup(name string, max int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	g := t.groups[name]
	if g == nil {
		g := newRunner(name, max)
		t.groups[name] = g
	} else {
		g.updateMax(max)
	}
}

func (t *Timer) AddJob(j *Job, group ...string) error {
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

	if len(group) > 0 {
		g := t.groups[group[0]]
		if g == nil {
			return fmt.Errorf("group named '%s' not exist", group[0])
		}
		j.js.runner = g
	} else {
		j.js.runner = dfRunner
	}

	err := t.parsingScheduleForJob(j)
	if err != nil {
		return fmt.Errorf("parsing schedule failed: %s", err)
	}

	j.timer = t
	t.jobs[j.name] = j
	t.queue.Push(j, j.sched.nextTicks())

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
		status     : StatusWaiting,
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
		status     : StatusWaiting,
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


