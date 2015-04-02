package parallel

import (
	"testing"
	"time"
)

func TestConditionVariable(t *testing.T) {
	cv := NewConditionVariable()
	ch := make(chan Empty)

	go func() {
		cv.Wait()
		ch <- Empty{}
	}()
	time.Sleep(time.Millisecond * 10)
	cv.NotifyOne()

	select {
	case <-ch:
		//t.Log("NotifyOne")
	case <-time.After(time.Millisecond * 10):
		t.Fatal("Wait has deadlock for NotifyOne")
	}

	for i := 0; i < 10; i++ {
		go func() {
			cv.Wait()
			ch <- Empty{}
		}()
	}
	time.Sleep(time.Millisecond * 10)
	cv.NotifyAll()

	time.Sleep(time.Millisecond * 10)
	cv.waiters = 1
	cv.NotifyAll()

	for i := 0; i < 10; i++ {
		select {
		case <-ch:
			//t.Log("NotifyAll", i)
		case <-time.After(time.Millisecond * 10):
			t.Fatal("Wait has deadlock for NotifyAll", i)
			break
		}
	}

	go func() {
		if cv.TimedWait(time.Millisecond * 10) {
			ch <- Empty{}
		}
	}()
	time.Sleep(time.Millisecond * 10)

	select {
	case <-ch:
		//t.Log("NotifyOne")
	case <-time.After(time.Millisecond * 20):
		t.Fatal("TimedWait has deadlock")
	}

	go func() {
		if !cv.TimedWait(time.Millisecond * 10) {
			ch <- Empty{}
		}
	}()
	time.Sleep(time.Millisecond * 5)
	cv.NotifyOne()

	select {
	case <-ch:
		//t.Log("NotifyOne")
	case <-time.After(time.Millisecond * 20):
		t.Fatal("TimedWait has deadlock")
	}
	cv.NotifyOne()
}
