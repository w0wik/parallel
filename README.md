# parallel

[![Build Status](https://travis-ci.org/w0wik/parallel.svg?branch=master)](https://travis-ci.org/w0wik/parallel)
[![GoDoc](https://godoc.org/github.com/w0wik/parallel?status.svg)](https://godoc.org/github.com/w0wik/parallel)

Collection of parallel and async wrappers for thread safe programming written in [golang](http://golang.org)

To install:

	go get github.com/w0wik/parallel

See docs for usage:
	https://godoc.org/github.com/w0wik/parallel

#### AsyncMap Example
```go
package main

import (
	"fmt"

	"github.com/w0wik/parallel"
)

func main() {
	m := parallel.NewAsyncMap(1)
	defer m.Close()
	m.Set(1) <- "hello"
	fmt.Println(<-m.Get(1))
	<-m.Delete(1)

}
```

#### ConditionVariable Example
```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/w0wik/parallel"
)

func wait(cv *parallel.ConditionVariable, n int) {
	cv.Wait()
	fmt.Println("Done waitng", n)
}

func main() {
	cv := parallel.NewConditionVariable()

	// 3 goroutines waits notify
	go wait(cv, 1)
	go wait(cv, 2)
	go wait(cv, 3)

	time.Sleep(time.Millisecond * 10)

	// 1 goroutine stops waiting
	cv.NotifyOne()

	// other goroutines stop waiting
	cv.NotifyAll()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
```

#### Monitor Example
```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/w0wik/parallel"
)

func main() {
	sl := make([]int, 0)
	mon := parallel.NewMonitor(&sl)

	go mon.Access(func(s *[]int) {
		*s = append(*s, 1)
	})
	go mon.Access(func(s *[]int) {
		*s = append(*s, 2)
	})
	time.Sleep(time.Millisecond * 15)
	mon.Access(func(s *[]int) {
		for _, v := range *s {
			fmt.Println(v)
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
```