package parallel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestSimpleReceiver struct {
	t   *testing.T
	fin chan int
}

func (tsr *TestSimpleReceiver) OnMessage(mess *Message) {
	assert.Equal(tsr.t, mess.Type, "hello", "Bad message type")
	tsr.fin <- 1
}

func TestSimpleMessage(t *testing.T) {
	hub := NewMessagesHub(100)
	defer hub.Stop()

	hub.SendMessage(nil, "hello")

	fin := make(chan int)
	hub.RegisterReceiver([]MessageType{"hello"}, &TestSimpleReceiver{t, fin})
	hub.SendMessage(nil, "hello")

	hub.SendMessage(nil, "bye")

	select {
	case <-fin:
	case <-time.After(time.Millisecond * 10):
		assert.Fail(t, "Message isn't received")
	}
}

type TestMultiReceiver struct {
	t   *testing.T
	typ string
	fin chan int
}

func (tmr *TestMultiReceiver) OnMessage(mess *Message) {
	assert.Equal(tmr.t, mess.Type, tmr.typ, "Bad message type")
	tmr.fin <- 1
}

func TestMultiMessages(t *testing.T) {
	fin := make(chan int)
	hub := NewMessagesHub(100)
	defer hub.Stop()
	hub.RegisterReceiver([]MessageType{"hello"}, &TestMultiReceiver{t, "hello", fin})
	hub.RegisterReceiver([]MessageType{"hello"}, &TestMultiReceiver{t, "hello", fin})
	hub.RegisterReceiver([]MessageType{"bye"}, &TestMultiReceiver{t, "bye", fin})
	hub.SendMessage(nil, "hello")
	hub.SendMessage(nil, "bye")

	for i := 0; i < 3; i++ {
		select {
		case <-fin:
		case <-time.After(time.Millisecond * 10):
			assert.Fail(t, "Message isn't received")
		}
	}
}

func TestDefautlHub(t *testing.T) {
	SendMessage(nil, "hello")
	defaultHub.Stop()
	defaultHub = nil
	RegisterReceiver([]MessageType{"hello"}, &TestMultiReceiver{t, "hello", nil})
	defaultHub.Stop()
	defaultHub = nil
	RegisterReceiveFunc([]MessageType{"hello"}, func(*Message) {})

}

type TestSenderReceiver struct {
	t   *testing.T
	typ string
	val int
	fin chan int
}

func (tmr *TestSenderReceiver) OnMessage(mess *Message) {
	assert.Equal(tmr.t, mess.Type, tmr.typ, "Bad message type")
	assert.Equal(tmr.t, mess.Sender, tmr.val, "Bad sender")
	tmr.fin <- 1
}

func TestMessageSenders(t *testing.T) {
	hub := NewMessagesHub(100)
	defer hub.Stop()

	fin := make(chan int)
	hub.RegisterReceiver([]MessageType{"hello"}, &TestSenderReceiver{t, "hello", 1, fin})
	hub.SendMessage(1, "hello")

	select {
	case <-fin:
	case <-time.After(time.Millisecond * 10):
		assert.Fail(t, "Message isn't received")
	}

}

func TestMessageReceiveFunc(t *testing.T) {
	hub := NewMessagesHub(100)
	defer hub.Stop()

	fin := make(chan int)
	hub.RegisterReceiveFunc([]MessageType{"hello"}, func(mess *Message) {
		assert.Equal(t, mess.Type, "hello", "Bad message type")
		fin <- 1
	})
	hub.SendMessage(nil, "hello")

	select {
	case <-fin:
	case <-time.After(time.Millisecond * 10):
		assert.Fail(t, "Message isn't received")
	}

}
