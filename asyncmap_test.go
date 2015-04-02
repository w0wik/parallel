package parallel

import (
	"fmt"
	"testing"
)

func TestAsyncMap(t *testing.T) {
	m := NewAsyncMap(1)
	m.Set(1) <- "hello"
	s := <-m.Get(1)
	if s != "hello" {
		t.Fatal("Not equal data")
	}
	if <-m.Len() != 1 {
		t.Fatal("Bad len")
	}

	<-m.Delete(1)

	if <-m.Len() != 0 {
		t.Fatal("Bad len after delete")
	}

	_, ok := <-m.Get(1)
	if ok {
		t.Fatal("Bad get from empty map")
	}

	for i := 0; i < 10; i++ {
		m.Set(i) <- i
	}
	listCount := 0
	for pair := range m.List() {
		listCount++
		if pair.First != pair.Second {
			t.Fatalf("%v!=%v", pair.First, pair.Second)
		}
	}
	if listCount != 10 {
		t.Fatal("List return few values", listCount)
	}

	m.Close()
	func() {
		defer func() {
			if res := recover(); res == nil {
				t.Fatal("no panic on closed map")
			}
		}()
		m.Get(1)
	}()
}

func BenchmarkSetToMap(b *testing.B) {
	m := NewAsyncMap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(i) <- 1
	}
}

func BenchmarkGetFromMap(b *testing.B) {
	m := NewAsyncMap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(i) <- 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-m.Get(i)
	}
}

func BenchmarkSetToMapParallel(b *testing.B) {
	m := NewAsyncMap(b.N)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(i) <- i
			i++
		}
	})
}

func BenchmarkGetFromMapParallel(b *testing.B) {
	m := NewAsyncMap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(i) <- 1
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-m.Get(i)
			i++
		}
	})
}

func ExampleAsyncMap() {
	m := NewAsyncMap(1)
	defer m.Close()
	m.Set(1) <- "hello"
	fmt.Println(<-m.Get(1))
	<-m.Delete(1)
}
