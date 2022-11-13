package etimer

import (
	"context"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/ziyht/eden_go/eerr"
)

type Job struct {
	name        string
	cb          JobFunc
	ctx         context.Context

	sched       schedule
	js          JobState

	timer       *Timer
}

// JobFunc is the timing called job function in timer.
type JobFunc = func(job *Job) error

type JobOpts struct {
	Name        string             // the name of the job, if not set, it will be replaced with internal generated uuid
	Ctx         context.Context    // context, to pass in needed parameters for the job
	Interval    time.Duration      // interval for job to run
	Pattern     string             // cron pattern to run job, has high priority than Interval
	CB          JobFunc            // callback function of job
	IsSingleton bool               // set singleton
	Times       int64              // set limit running times
	status      int32
}

// createJob creates and adds a timing job to the timer.
func createJob(in JobOpts) (*Job) {
	j := &Job{
		name:        in.Name,
		cb:          in.CB,
		ctx:         in.Ctx,
	}

	if j.name == "" {
		j.name = uuid.NewV1().String()
	}

	j.js.errs = newErrInfos()
	j.js.name = j.name
	j.js.pattern = strings.TrimSpace(in.Pattern)
	j.js.setStatus(in.status)
	j.js.setInterval(in.Interval)
	j.js.setSingleton(in.IsSingleton)
	if in.Times > 0 {
		j.js.setTimes(in.Times)
	}

	if j.js.pattern != "" {
		j.js.setInterval(time.Second)
	}

	return j
}

// Errorf - format a err and record it to the JobState
func (j *Job) Errorf(message string, args ...interface{}) {
	j.js.recordError(eerr.NewSkipf(1, message, args...))
}

// RecordError - record the err to the JobState
func (j *Job) RecordError(err error) {
	j.js.recordError(eerr.PackSkip(1, err))
}

// Start starts the job.
func (j *Job) Start() {
	j.js.setStatus(StatusReady)
}

// Stop stops the job.
func (j *Job) Stop() {
	j.js.setStatus(StatusStopped)
}

func (j *Job) Close(){
	j.js.setStatus(StatusClosed)
}

func (j *Job) IsSingleton() bool {
	return j.js.IsSingleton()
}

func (j *Job) SetSingleton(enable bool) {
	j.js.setSingleton(enable)
}

func (j *Job) Ctx() context.Context {
	return j.ctx
}

func (j *Job) SetTimes(times int64) {
	j.js.setTimes(times)
}

func (j *Job) State() *JobState {
	return &j.js
}

func (j *Job) exec_once() {
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
}
