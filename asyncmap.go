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

// Async map provides asynchronous and thread-safe access to map.
// All funcs of this class return channels to interact with map
type AsyncMap struct {
	cache    map[interface{}]interface{}
	commands chan *async_command
	closed   bool
}

// NewAsyncMap creates new async map.
// cache indicates a size of channel of interim values.
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

// Set returns channel to set the value
func (am *AsyncMap) Set(key interface{}) chan<- interface{} {
	am.check_closed()
	ch := make(chan interface{}, 1)
	am.commands <- &async_command{ch, async_command_set, key}
	return ch
}

// Get return channel to get the value
// If channel closed then map don't have value by the key.
func (am *AsyncMap) Get(key interface{}) <-chan interface{} {
	am.check_closed()
	ch := make(chan interface{}, 1)
	am.commands <- &async_command{ch, async_command_get, key}
	return ch
}

// Delete ask to delete the element and returns channel which indicates that the element is deleted
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

// Len returns channel to get the len of map
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

// Close closes internal channels and finish work goroutine.
func (am *AsyncMap) Close() {
	am.closed = true
	close(am.commands)
}
