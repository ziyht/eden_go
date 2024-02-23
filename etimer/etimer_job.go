package etimer

import (
	"context"
	"time"

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
	j.js.setStatus(StatusWaiting)
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
