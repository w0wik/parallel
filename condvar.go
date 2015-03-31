package parallel

import (
	"sync/atomic"
	"time"
)

type Empty struct{}
type ConditionVariable struct {
	ch      chan Empty
	waiters int32
}

func NewConditionVariable() *ConditionVariable {
	ret := new(ConditionVariable)
	ret.ch = make(chan Empty)
	return ret
}

func (s *ConditionVariable) NotifyOne() {
	select {
	case <-s.ch:
		atomic.AddInt32(&s.waiters, -1)
	default:
	}
}

func (s *ConditionVariable) NotifyAll() {
	wt := int(atomic.LoadInt32(&s.waiters))
	for i := 0; i < wt; i++ {
		select {
		case <-s.ch:
			atomic.AddInt32(&s.waiters, -1)
		default:
			atomic.StoreInt32(&s.waiters, 0)
			break
		}
	}
}

func (s *ConditionVariable) Wait() {
	atomic.AddInt32(&s.waiters, 1)
	s.ch <- Empty{}
}

func (s *ConditionVariable) TimedWait(tm time.Duration) (is_timeout bool) {
	atomic.AddInt32(&s.waiters, 1)
	select {
	case s.ch <- Empty{}:
		is_timeout = false
	case <-time.After(tm):
		is_timeout = true
	}
	return
}
