package parallel

import (
	"log"
	"reflect"
	"sync"
)

// Monitor provides exclusive access to object.
type Monitor struct {
	obj reflect.Value
	mut sync.Mutex
}

// NewMonitor creates new monitor for obj.
func NewMonitor(obj interface{}) *Monitor {
	ret := new(Monitor)
	ret.obj = reflect.ValueOf(obj)
	return ret
}

// Access locks object and run access_fun with monitored object as argument.
func (m *Monitor) Access(access_fun interface{}) {

	fn := reflect.ValueOf(access_fun)

	if fn.Kind() != reflect.Func {
		log.Panicf("Argument of Access must be func %v ", fn)
		return
	}

	m.mut.Lock()
	defer m.mut.Unlock()
	fn.Call([]reflect.Value{m.obj})
}
