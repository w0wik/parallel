package parallel

import (
	"sync/atomic"
	"time"
)

// Empty is empty struct for channels
type Empty struct{}

// ConditionVariable provides a mechanism for notifies from one goroutine to another
type ConditionVariable struct {
	ch      chan Empty
	waiters int32
}

// NewConditionVariable creates new condition variable
func NewConditionVariable() *ConditionVariable {
	ret := new(ConditionVariable)
	ret.ch = make(chan Empty)
	return ret
}

// NotifyOne notifies to one waiting goroutine
func (s *ConditionVariable) NotifyOne() {
	select {
	case <-s.ch:
		atomic.AddInt32(&s.waiters, -1)
	default:
	}
}

// NotifyAll notfies to all waiting goroutines
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

// Wait waits for notify from other goroutine
func (s *ConditionVariable) Wait() {
	atomic.AddInt32(&s.waiters, 1)
	s.ch <- Empty{}
}

// TimedWait waits for notify from other goroutine. It ends waiting by timeout tm.
// If TimedWait ends waiting by timeout then it return true
func (s *ConditionVariable) TimedWait(tm time.Duration) (isTimeout bool) {
	atomic.AddInt32(&s.waiters, 1)
	select {
	case s.ch <- Empty{}:
		isTimeout = false
	case <-time.After(tm):
		isTimeout = true
	}
	return
}
