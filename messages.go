package parallel

import "sync"

type MessageType interface{}
type MessageArg interface{}

type Message struct {
	Type   MessageType
	Sender MessageSender
	Args   []MessageArg
}
type MessageReceiver interface {
	OnMessage(*Message)
}
type MessageSender interface {
	//	SendMessage(*Message)
}
type MessagesHub struct {
	receivers *AsyncMap
	messages  chan *Message
}

func NewMessagesHub(cache int) *MessagesHub {
	ret := new(MessagesHub)
	ret.receivers = NewAsyncMap(cache)
	ret.messages = make(chan *Message, cache)

	go ret.doMessages()
	return ret
}

func (mh *MessagesHub) Stop() {
	close(mh.messages)
}

func (mh *MessagesHub) doMessages() {
	for mess := range mh.messages {
		recvrs, ok := <-mh.receivers.Get(mess.Type)
		if ok {
			for _, r := range recvrs.([]MessageReceiver) {
				loc_mess := *mess // copy message for security reasons
				go r.OnMessage(&loc_mess)
			}
		}
	}
	mh.receivers.Close()
}

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

func (mh *MessagesHub) SendMessage(sender MessageSender, typ MessageType, args ...MessageArg) {
	mh.messages <- &Message{typ, sender, args}
}

/*func (mh *MessagesHub) RegisterSender(sender MessageSender) error {
	return nil
}*/

var (
	defaultHub    *MessagesHub
	defaultHubMut sync.Mutex
)

func RegisterReceiver(types []MessageType, receiver MessageReceiver) {
	defaultHubMut.Lock()
	if defaultHub == nil {
		defaultHub = NewMessagesHub(100)
	}
	defaultHubMut.Unlock()

	defaultHub.RegisterReceiver(types, receiver)
}

func SendMessage(sender MessageSender, typ MessageType, args ...MessageArg) {
	defaultHubMut.Lock()
	if defaultHub == nil {
		defaultHub = NewMessagesHub(100)
	}
	defaultHubMut.Unlock()
	defaultHub.SendMessage(sender, typ, args...)
}
