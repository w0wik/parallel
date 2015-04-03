package parallel

import (
	"fmt"
	"reflect"
	"sync"
)

// Monitor provides exclusive access to object.
type Monitor struct {
	obj reflect.Value
	mut sync.RWMutex
}

// NewMonitor creates new monitor for obj.
func NewMonitor(obj interface{}) *Monitor {
	ret := new(Monitor)
	ret.obj = reflect.ValueOf(obj)
	return ret
}

// Access locks object and run accessFun with monitored object as argument.
func (m *Monitor) Access(accessFun interface{}) {
	m.access(false, accessFun)
}

// RAccess locks object for read and run accessFun with monitored object as argument.
func (m *Monitor) RAccess(accessFun interface{}) {
	m.access(true, accessFun)
}

func (m *Monitor) access(readLock bool, accessFun interface{}) {

	fn := reflect.ValueOf(accessFun)

	if fn.Kind() != reflect.Func {
		panic(fmt.Sprintf("Argument of Access must be func %v ", fn))
	}

	if readLock {
		m.mut.RLock()
		defer m.mut.RUnlock()
	} else {
		m.mut.Lock()
		defer m.mut.Unlock()
	}
	fn.Call([]reflect.Value{m.obj})
}
