package etimer

import (
	"strconv"
	"time"
)


const (
	StatusReady                            = 0                    // Job or Timer is ready for running.
	StatusRunning                          = 1                    // Job or Timer is already running.
	StatusStopped                          = 2                    // Job or Timer is stopped.
	StatusClosed                           = -1                   // Job or Timer is closed and waiting to be deleted.
	panicExit                internalPanic = "exit"               // panicExit is used for custom job exit with panic.
	defaultTimerInterval                   = "100"                // defaultTimerInterval is the default timer interval in milliseconds.
	commandEnvKeyForInterval               = "gf.gtimer.interval" // commandEnvKeyForInterval is the key for command argument or environment configuring default interval duration for timer.
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
