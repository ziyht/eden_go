package etimer

import (
	"sync/atomic"
	"time"
)

type JobState struct {
	status      int32

	runnings    uint64         // 运行完成的次数
	failures    uint64         // 失败次数
	successs    uint64         // 成功次数

	nextStart   time.Time      // 下一次开始时间
	lastStart   time.Time      // 上一次开始时间
	lastOver    time.Time      // 上一次结束时间
	lastCost    time.Duration  // 上一次运行耗时
	lastSuccess time.Time      // 上一次成功的时间
	lastFailed  time.Time      // 上一次失败的时间
	lastError   error          // 上一次失败的错误信息
}

func (js *JobState)setStatus(s int32) { atomic.StoreInt32(&js.status, s) }
func (js *JobState)addFailure() {  atomic.AddUint64(&js.runnings, 1); atomic.AddUint64(&js.failures, 1)}
func (js *JobState)addSuccess() {  atomic.AddUint64(&js.runnings, 1); atomic.AddUint64(&js.successs, 1)}

func (js *JobState)setStart(t time.Time) { js.lastStart = t }


func (js *JobState)Status() int { return int(atomic.LoadInt32(&js.status)) }

func (js *JobState)Runnings() uint64 { return atomic.LoadUint64(&js.runnings) }
func (js *JobState)Failures() uint64 { return atomic.LoadUint64(&js.failures) }
func (js *JobState)Successs() uint64 { return atomic.LoadUint64(&js.successs) }

func (js *JobState)NextStart() time.Time { return js.nextStart }
