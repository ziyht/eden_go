package etimer

import (
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/ziyht/eden_go/eerr"
)

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

	j.js.runner = dfRunner
	j.js.errs   = newErrInfos()
	j.js.name   = j.name
	j.js.pattern = strings.TrimSpace(in.Pattern)
	j.js.setStatus(in.status)
	j.js.setInterval(in.Interval)
	j.js.setSingleton(in.IsSingleton)
	if in.Times > 0 {
		j.js.setTimes(in.Times)
	}

	if j.js.pattern != "" {
		j.js.setInterval(time.Second)
		// if in.Interval == 0 {
		// 	j.js.setInterval(10 * time.Second)
		// } else if in.Interval < time.Second{
		// 	j.js.setInterval(time.Second)
		// } else {
		// 	j.js.setInterval(in.Interval)
		// }
	}

	return j
}

func (j *Job) exec_once() {
	if j.js.checkLimit(-1){
		return
	}
	
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
			j.js.setStatus(StatusWaiting)
		}
	}()

	err   := j.cb(j)
	j.js.addRunningOver(start, time.Now(), err)
}