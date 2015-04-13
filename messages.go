package parallel

import "sync"

// MessageType provides message type.
type MessageType interface{}

// MessageArg provides message arg.
type MessageArg interface{}

// Message is a object thats is used for send/receive message.
type Message struct {
	Type   MessageType
	Sender MessageSender
	Args   []MessageArg
}

// MessageReceiver is the interface implemented by an object that can receive messages.
type MessageReceiver interface {
	OnMessage(*Message)
}

// MessageSender is the interface implemented by an object that can send messages.
type MessageSender interface {
}

// MessagesHub manages messages.
type MessagesHub struct {
	receivers *AsyncMap
	messages  chan *Message
}

// NewMessagesHub creates new MessagesHub. cache provides a size of worker channels.
func NewMessagesHub(cache int) *MessagesHub {
	ret := new(MessagesHub)
	ret.receivers = NewAsyncMap(cache)
	ret.messages = make(chan *Message, cache)

	go ret.doMessages()
	return ret
}

// Stop stops all worker goroutines and closes worker channels.
func (mh *MessagesHub) Stop() {
	close(mh.messages)
}

func (mh *MessagesHub) doMessages() {
	for mess := range mh.messages {
		recvrs, ok := <-mh.receivers.Get(mess.Type)
		if ok {
			for _, r := range recvrs.([]MessageReceiver) {
				locMess := *mess // copy message for security reasons
				go r.OnMessage(&locMess)
			}
		}
	}
	mh.receivers.Close()
}

// RegisterReceiver registers new MessageReceiver.
func (mh *MessagesHub) RegisterReceiver(types []MessageType, receiver MessageReceiver) {
	for _, typ := range types {
		receivers := <-mh.receivers.Get(typ)
		if receivers != nil {
			mh.receivers.Set(typ) <- append(receivers.([]MessageReceiver), receiver)
		} else {
			mh.receivers.Set(typ) <- []MessageReceiver{receiver}
		}
	}
}

type oneFuncReseiver func(*Message)

func (f oneFuncReseiver) OnMessage(m *Message) {
	f(m)
}

// RegisterReceiveFunc register single function for receive Messages.
func (mh *MessagesHub) RegisterReceiveFunc(types []MessageType, receiver func(*Message)) {
	mh.RegisterReceiver(types, oneFuncReseiver(receiver))
}

// SendMessage sends message to registered MessageReceivers.
func (mh *MessagesHub) SendMessage(sender MessageSender, typ MessageType, args ...MessageArg) {
	mh.messages <- &Message{typ, sender, args}
}

var (
	defaultHub    *MessagesHub
	defaultHubMut sync.Mutex
)

// RegisterReceiver registers new MessageReceiver.
func RegisterReceiver(types []MessageType, receiver MessageReceiver) {
	defaultHubMut.Lock()
	if defaultHub == nil {
		defaultHub = NewMessagesHub(100)
	}
	defaultHubMut.Unlock()

	defaultHub.RegisterReceiver(types, receiver)
}

// RegisterReceiveFunc register single function for receive Messages.
func RegisterReceiveFunc(types []MessageType, receiver func(*Message)) {
	RegisterReceiver(types, oneFuncReseiver(receiver))
}

// SendMessage sends message to registered MessageReceivers.
func SendMessage(sender MessageSender, typ MessageType, args ...MessageArg) {
	defaultHubMut.Lock()
	if defaultHub == nil {
		defaultHub = NewMessagesHub(100)
	}
	defaultHubMut.Unlock()
	defaultHub.SendMessage(sender, typ, args...)
}
