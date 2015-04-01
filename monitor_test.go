package parallel

import (
	"fmt"
	"sync"
	"testing"
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
	if len(obj) < 10 {
		t.Fatal("Failed")
	}
	func() {
		defer func() {
			if res := recover(); res == nil {
				t.Fatal("no panic on bad access map")
			}
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
