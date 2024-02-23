package etimer

import (
	"fmt"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/panjf2000/ants/v2"
)

type runner struct {
  name_  string
	rtp    *ants.Pool

	sigs   chan int

	tasksd map[*Job]bool
	tasksl *singlylinkedlist.List

	mu     sync.Mutex
}

var dfRunner = newRunner("(INNER_UNLIMITED)", -1)

func newRunner(name string, max int) *runner {
	pool, _ := ants.NewPool(max)

	r := &runner{
	  name_: name,
		rtp  : pool,
		sigs : make(chan int, 10),
		tasksd: map[*Job]bool{},
		tasksl: singlylinkedlist.New(),
	}

	go r.loop()

	return r
}

func (r *runner) _put(job *Job) bool {
	prev := r.tasksd[job]
	if prev {
		return false
	}

	r.tasksd[job] = true
	r.tasksl.Append(job)

	return true
}

func (r *runner) _removeFirst() *Job {
	a, e := r.tasksl.Get(0)
	if e {
		j := a.(*Job)
		delete(r.tasksd, j)
		r.tasksl.Remove(0)

		return j
	}

	return nil
}

func (r *runner) updateMax(max int) {
	r.rtp.Tune(max)
	r.notifyWakeup()
}

// func (r *runner) name() string {
// 	return r.name_
// }

func (r *runner) submit(job *Job){
	r.mu.Lock()
	defer r.mu.Unlock()

	if job.IsSingleton(){
		p, r_ := job.js.getPendRun()
		if p == 0 {
			job.js.addPendRun(1, 0)
		}
		if r_ == 0 {
			r._put(job)
		}
	} else {
		r._put(job)
		job.js.addPendRun(1, 0)
	}

	r.notifyWakeup()
}

func (r *runner) notifyWakeup() {
  select {
    case r.sigs <- 1:
		  return
    default: return
  }
}

// func (r *runner) notifyQuit() {
// 	r.sigs <- 0
// }

func (r *runner) running() int {
	return r.rtp.Running()
}

func (r *runner) pending() int {
	return r.rtp.Waiting()
}

func (r *runner) max() int {
	return r.rtp.Cap()
}

func (r *runner)String() string {
	return fmt.Sprintf("r:%d p:%d m:%d", r.running(), r.pending(), r.max())
}

func (r *runner)_waitSig() int {

	sig := <- r.sigs
	if sig == 0 { 
		return 0 
	}

	// here remove all sigs
	for {
		select {
			case sig = <- r.sigs: 
				if sig == 0 {
					return 0
				}
			default: return sig
		}
	}
}

func (r *runner)loop() {
	for {
		sig := r._waitSig()

		// quit?
		if sig == 0 {
			goto quit
		}

		if r.rtp.Free() == 0 {
			continue
		}

		if r.tasksl.Size() == 0 {
			continue
		}

		r.mu.Lock()
		job := r._removeFirst()

		p_, r_ := job.js.getPendRun()
		if p_ == 0 || r_ == 1 {
			r.mu.Unlock()
			continue
		}

		job.js.addPendRun(-1, 1)
		r.mu.Unlock()

		err := r.rtp.Submit(func (){
			job.exec_once()

			r.mu.Lock()
			p_, _ := job.js.addPendRun(0, -1)
			if p_ > 0 {
				r._put(job)
			}
			r.mu.Unlock()

			r.notifyWakeup()
		})

		if err != nil {
			job.js.errs.setError(err, time.Now())

			r.mu.Lock()
			p_, _ := job.js.addPendRun(0, -1)
			if p_ > 0 {
				r._put(job)
			}
			r.mu.Unlock()
		}
	}

quit:
}

