package parallel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsyncMap(t *testing.T) {
	m := NewAsyncMap(1)
	m.Set(1) <- "hello"
	s := <-m.Get(1)
	assert.Equal(t, s, "hello", "Get() returning bad data from map")

	assert.Equal(t, <-m.Len(), 1, "Len() returning uncorrect value")

	<-m.Delete(1)

	assert.Equal(t, <-m.Len(), 0, "Len() returning uncorrect value after delete")

	_, ok := <-m.Get(1)
	assert.False(t, ok, "Get() returning data from empty map")

	for i := 0; i < 10; i++ {
		m.Set(i) <- i
	}
	listCount := 0
	for pair := range m.List() {
		listCount++
		assert.Equal(t, pair.First, pair.Second, "List() returning bad data")
	}
	assert.Equal(t, listCount, 10, "List() returning few values")

	m.Close()
	func() {
		defer func() {
			assert.NotNil(t, recover(), "No panic on closed map")
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
