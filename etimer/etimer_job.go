package etimer

import (
	"context"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
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
