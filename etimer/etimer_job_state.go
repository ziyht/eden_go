package etimer

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type JobState struct {
	mu          sync.Mutex
	runner      *runner

	name        string         // the name of the job 
	run_pend    uint64

	status      int32
	singleton   int32
	interval    time.Duration
	pattern     string
	times       int64          // Limit running times.
	leftTimes   int64          // Left running times, init by times
	
	runnings    uint64         // 运行的次数
	failures    uint64         // 失败次数
	successs    uint64         // 成功次数

	nextStart   time.Time      // 下一次开始时间
	lastStart   time.Time      // 上一次开始时间
	lastEnd     time.Time      // 上一次结束时间
	lastCost    time.Duration  // 上一次运行耗时
	lastSuccess time.Time      // 上一次成功的时间
	lastFailure time.Time      // 上一次失败的时间
	lastError   error          // 上一次失败的错误信息

	errs        *errInfos
}

func (js *JobState)setStatus(s int32) { atomic.StoreInt32(&js.status, s) }
func (js *JobState)setInterval(du time.Duration) { js.interval = du }
func (js *JobState)setStatusCas(s1, s2 int32) bool { return atomic.CompareAndSwapInt32(&js.status, s1, s2) }
func (js *JobState)setSingleton(enable bool) {
	if enable {
		atomic.StoreInt32(&js.singleton, 1)
	} else {
		atomic.StoreInt32(&js.singleton, 0)
	}
}
func (js *JobState)setTimes(times int64) {
	if times > 1 {
		atomic.StoreInt64(&js.times, times)
		atomic.StoreInt64(&js.leftTimes, times)
	}
}
func (js *JobState)getPendRun()(int32, int32){
	val := atomic.LoadUint64(&js.run_pend)
	return int32(val >> 32), int32(val)
}
func (js *JobState)addPendRun(p, r int32) (int32, int32) {
	var pending, running int32

	for {
		val := atomic.LoadUint64(&js.run_pend)
		pending = int32(val >> 32)
		running = int32(val)

		if p == 0 && r == 0 {
			break
		}

		pending += p
		running += r
		
		new_val := uint64(pending) << 32 | uint64(running)

		if atomic.CompareAndSwapUint64(&js.run_pend, val, new_val) {
			break
		}
	}
	
	return pending, running
}
func (js *JobState)addRunning(start time.Time) {
	atomic.AddUint64(&js.runnings, 1)
	js.lastStart = start
}

func (js *JobState)checkLimit(cnt int)bool{
	if js.times <= 0 {
		return false
	}

	left := atomic.AddInt64(&js.leftTimes, int64(cnt))
	if left > 0 {
		return false
	}

	if left == 0 {
		js.setStatus(StatusStopped)
		return false
	}

	js.setStatus(StatusStopped)
	return true
}

func (js *JobState)addRunningOver(start, end time.Time, err error) {
	js.lastStart = start
	js.lastEnd   = end
	js.lastCost  = end.Sub(start)

	if err == nil {
		atomic.AddUint64(&js.successs, 1)
		js.lastSuccess = end
	} else {
		atomic.AddUint64(&js.failures, 1)
		js.lastFailure = end
	}
}

func (js *JobState)recordError(err error, t ...time.Time) {
	if len(t) > 0 {
		js.errs.setError(err, t[0])
	} else {
		js.errs.setError(err, time.Now())
	}

	js.lastError = err
}


func (js *JobState)Name()      string { return js.name }
func (js *JobState)Interval()  time.Duration { return js.interval }
func (js *JobState)Status()    int    { return int(atomic.LoadInt32(&js.status)) }
func (js *JobState)StatusStr() string { return statusString(atomic.LoadInt32(&js.status)) }
func (js *JobState)Runnings()  uint64 { return atomic.LoadUint64(&js.runnings) }
func (js *JobState)Failures()  uint64 { return atomic.LoadUint64(&js.failures) }
func (js *JobState)Successs()  uint64 { return atomic.LoadUint64(&js.successs) }

func (js *JobState)Times()       int64 { return atomic.LoadInt64(&js.times) }  // return limit times of job
func (js *JobState)LeftTimes()   int64 { return atomic.LoadInt64(&js.leftTimes) }
func (js *JobState)IsSingleton() bool  { return atomic.LoadInt32(&js.singleton) > 0 }

func (js *JobState)NextStart() time.Time { return js.nextStart }

func (js *JobState)LastStart() time.Time { return js.lastStart }
func (js *JobState)LastEnd()   time.Time { return js.lastEnd }
func (js *JobState)LastCost()  time.Duration { return js.lastCost }

func (js *JobState)LastSuccess()  time.Time { return js.lastSuccess }
func (js *JobState)LastFailure()  time.Time { return js.lastFailure }
func (js *JobState)LastError(clear ...bool) error { err := js.lastError; if len(clear) > 0 && clear[0] { js.lastError = nil }; return err}

func (js *JobState)Errors() []*errInfo { return js.errs.Errs() }

func (js *JobState)Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			io.WriteString(s, js.formatLevel1())
		case s.Flag('+'):
			if verb == 's' {
				io.WriteString(s, js.formatLevel2())
			} else {
				io.WriteString(s, js.formatLevel2())
			}
		default:
			io.WriteString(s, js.formatLevel1())
		}
	}
}

const jsTimeFormat = "2006-01-02T15:04:05.999999"

func (js *JobState)formatLevel1() string {
	return fmt.Sprintf("lastStart: %s, lastEnd: %-s, next: %-s, runnings: %d(%d|%d)", 
		js.lastStart.Format(jsTimeFormat),
		js.lastEnd.Format(jsTimeFormat),
		js.nextStart.Format(jsTimeFormat),
    js.runnings, js.successs, js.failures)
}

func (js *JobState)formatLevel2() string {
	return fmt.Sprintf(`name       : %s
status     : %s
runnings   : %d
failures   : %d
successs   : %d
nextStart  : %s
lastStart  : %s
lastEnd    : %s
lastCost   : %s
lastSuccess: %s
lastFailure: %s
lastError  : %s
`, 
		js.name,
		statusString(js.status),
		js.runnings,
		js.failures,
		js.successs,
		js.nextStart.Format(jsTimeFormat),
		js.lastStart.Format(jsTimeFormat),
		js.lastEnd.Format(jsTimeFormat),
		js.lastCost,
		js.lastSuccess.Format(jsTimeFormat),
		js.lastFailure.Format(jsTimeFormat),
    js.lastError)
}