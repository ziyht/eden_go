package etimer

import (
	"context"
	"sync/atomic"
)

type Job struct {
	name        string
	cb          JobFunc
	ctx         context.Context
	ticks       int64           // The job runs every tick.
	nextTicks   int64           // Next run ticks of the job.
	times       int64           // Limit running times.
	leftTimes   int64           // Left running times, init by times
	singleton   int32
	state       JobState

	timer       *Timer
}

// JobFunc is the timing called job function in timer.
type JobFunc = func(ctx context.Context) error

func (j *Job) Run() {
	if j.times >= 0 {
		leftRunTimes := atomic.AddInt64(&j.leftTimes, -1)
		if leftRunTimes < 0 {
			j.state.setStatus(StatusClosed)
			return
		}
	}

	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if exception != panicExit {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						panic(v)
					} else {
						panic(gerror.Newf(`exception recovered: %+v`, exception))
					}
				} else {
					j.Close()
					return
				}
			}
			if j.Status() == StatusRunning {
				j.state.setStatus(StatusReady)
			}
		}()
		j.cb(j.ctx)
	}()
}

// doCheckAndRunByTicks checks the if job can run in given timer ticks,
// it runs asynchronously if the given `currentTimerTicks` meets or else
// it increments its ticks and waits for next running check.
func (j *Job) doCheckAndRunByTicks(currentTimerTicks int64) {
	// Ticks check.
	if currentTimerTicks < atomic.LoadInt64(&j.nextTicks) {
		return
	}
	atomic.StoreInt64(&j.nextTicks, currentTimerTicks + j.ticks)
	// Perform job checking.
	switch j.state.Status() {
	  case StatusRunning: if j.IsSingleton() { return }
	  case StatusReady  : if !atomic.CompareAndSwapInt32(&j.state.status, StatusReady, StatusRunning){ return }
	  case StatusStopped: return
	  case StatusClosed : return
	}
	// Perform job running.
	j.Run()
}

func (j *Job) Status()int{
	return j.state.Status()
}

// Start starts the job.
func (j *Job) Start() {
	j.state.setStatus(StatusReady)
}

// Stop stops the job.
func (j *Job) Stop() {
	j.state.setStatus(StatusStopped)
}

func (j *Job) Close(){
	j.state.setStatus(StatusClosed)
}

func (j *Job) IsSingleton() bool {
	return atomic.LoadInt32(&j.singleton) > 0
}

func (j *Job) SetSingleton(enable bool) {
	if enable {
		atomic.StoreInt32(&j.singleton, 1)
	} else {
		atomic.StoreInt32(&j.singleton, 0)
	}
}

func (j *Job) Ctx() context.Context {
	return j.ctx
}

func (j *Job) SetTimes(times int64) {
	atomic.StoreInt64(&j.times, times)
	atomic.StoreInt64(&j.leftTimes, times)
}

func (j *Job) State() *JobState {
	return &j.state
}
