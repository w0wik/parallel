package parallel

const (
	async_command_get int = iota
	async_command_set
	async_command_delete
	async_command_len
)

type async_command struct {
	ch  chan interface{}
	typ int
	key interface{}
}

type AsyncMap struct {
	cache    map[interface{}]interface{}
	commands chan *async_command
	closed   bool
}

func NewAsyncMap(cache int) *AsyncMap {
	ret := new(AsyncMap)
	ret.cache = make(map[interface{}]interface{})
	ret.commands = make(chan *async_command, cache)
	go func() {
		for cmd := range ret.commands {
			switch cmd.typ {
			case async_command_get:
				if value, ok := ret.cache[cmd.key]; ok {
					cmd.ch <- value
				} else {
					close(cmd.ch)
				}
			case async_command_set:
				ret.cache[cmd.key] = <-cmd.ch
			case async_command_delete:
				delete(ret.cache, cmd.key)
				cmd.ch <- Empty{}
			case async_command_len:
				cmd.ch <- len(ret.cache)
			}
		}
		ret.cache = nil
	}()
	return ret
}

func (am *AsyncMap) check_closed() {
	if am.closed {
		panic("Async map is closed")
	}
}

func (am *AsyncMap) Set(key interface{}) chan<- interface{} {
	am.check_closed()
	ch := make(chan interface{}, 1)
	am.commands <- &async_command{ch, async_command_set, key}
	return ch
}

func (am *AsyncMap) Get(key interface{}) <-chan interface{} {
	am.check_closed()
	ch := make(chan interface{}, 1)
	am.commands <- &async_command{ch, async_command_get, key}
	return ch
}

func (am *AsyncMap) Delete(key interface{}) <-chan Empty {
	am.check_closed()
	ich := make(chan interface{}, 1)
	ch := make(chan Empty, 1)
	go func() {
		ch <- (<-ich).(Empty)
	}()

	am.commands <- &async_command{ich, async_command_delete, key}
	return ch
}

func (am *AsyncMap) Len() <-chan int {
	am.check_closed()
	ich := make(chan interface{}, 1)
	ch := make(chan int, 1)
	go func() {
		ch <- (<-ich).(int)
	}()
	am.commands <- &async_command{ich, async_command_len, 0}
	return ch
}

func (am *AsyncMap) Close() {
	am.closed = true
	close(am.commands)
}
