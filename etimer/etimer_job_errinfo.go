package etimer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/ziyht/eden_go/eerr"
)

type errInfo struct {
	first time.Time
	last  time.Time
	cnt   uint64
	frame eerr.Frame
	err   error
}

func (e *errInfo) addCount(v uint64) { atomic.AddUint64(&e.cnt, v) }

func (e *errInfo) First() time.Time { return e.first }
func (e *errInfo) Last()  time.Time { return e.last }
func (e *errInfo) Count() uint64    { return e.cnt }
func (e *errInfo) File()  string    { return e.frame.File() }
func (e *errInfo) Func()  string    { return e.frame.Func() }
func (e *errInfo) Line()  int       { return e.frame.Line() }

type errInfos struct {
	mu    sync.Mutex
	errs  map[uint64]*errInfo
	errsl []*errInfo
}

func newErrInfos() *errInfos {
	return &errInfos{errs: map[uint64]*errInfo{},}
}

func (es *errInfos)setError(err error, t time.Time) {
	if err == nil {
		return 
	}

	es.mu.Lock()
	defer es.mu.Unlock()

	id := eerr.Id(err)
	fd := es.errs[id]
	if fd != nil {
		fd.addCount(1)
		fd.last = t
		return
	}

	newEI := &errInfo{
		first: t,
		last : t,
		cnt  : 1,
		err  : err,
	}

	e, ok := err.(eerr.Error) 
	if ok {
		newEI.frame = e.StaskCause()
		newEI.err   = e.Unpack()
	} 

	es.errs[id] = newEI
	es.errsl = append(es.errsl, newEI)

	return 
}

func (es *errInfos)Errs() []*errInfo {
	return es.errsl
}