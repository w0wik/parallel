# parallel

[![Build Status](https://travis-ci.org/w0wik/parallel.svg?branch=master)](https://travis-ci.org/w0wik/parallel)

Collection of parallel and async wrappers for thread safe programming written in golang (http://golang.org)

To install:

	go get github.com/w0wik/parallel


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

