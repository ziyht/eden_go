package etimer

import (
	"context"
	"strconv"
	"time"
)


const (
	StatusWaiting                          = 0                    // Job or Timer is waiting for next schedule to run.
	StatusRunning                          = 1                    // Job or Timer is already running.
	StatusStopped                          = 2                    // Job or Timer is stopped.
	StatusPending                          = 3                    // Job is submitted and pending to running.
	StatusClosed                           = -1                   // Job or Timer is closed and waiting to be deleted.
	panicExit                internalPanic = "exit"               // panicExit is used for custom job exit with panic.
	defaultTimerInterval                   = "100"                // defaultTimerInterval is the default timer interval in milliseconds.
)

var (
	defaultInterval = getDefaultInterval()
	defaultTimer    = newTimer()
)

func getDefaultInterval() time.Duration {
	n, _ := strconv.Atoi(defaultTimerInterval)
	return time.Duration(n) * time.Millisecond
}

// DefaultOptions creates and returns a default options object for Timer creation.
func DefaultOptions() TimerOptions {
	return TimerOptions{
		Interval: defaultInterval,
	}
}

func NewTimer(opt ... TimerOptions)*Timer {
	return newTimer(opt...)
}

func NewJob(opt *JobOpts)*Job {
	return createJob(*opt)
}

func AddJob(j *Job) error {
	return defaultTimer.AddJob(j)
}

func AddInterval(ctx context.Context, interval time.Duration, cb JobFunc, singleton ...bool) *Job {
	return defaultTimer.AddInterval(ctx, interval, cb, singleton...)
}

func AddCron(ctx context.Context, pattern string, cb JobFunc, singleton ...bool) (*Job, error) {
	return defaultTimer.AddCron(ctx, pattern, cb, singleton...)
}