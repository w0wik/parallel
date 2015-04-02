package parallel

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	var obj []int
	m := NewMonitor(&obj)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			m.Access(func(s *[]int) {
				*s = append(*s, i)
			})
		}(i)
	}
	wg.Wait()
	assert.False(t, len(obj) < 10, "Error of append to monitored slice")

	func() {
		defer func() {
			assert.NotNil(t, recover(), "No panic on bad access func")
		}()

		m.Access(nil)
	}()
}

func ExampleMonitor() {
	obj := []int{1, 2, 3}
	m := NewMonitor(obj)
	go m.Access(func(s []int) {
		fmt.Println(s)
	})
}
