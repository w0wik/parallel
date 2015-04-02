package parallel

const (
	asyncCommandGet int = iota
	asyncCommandSet
	asyncCommandDelete
	asyncCommandLen
	asyncCommandList
)

type asyncCommand struct {
	ch  chan interface{}
	typ int
	key interface{}
}

// Pair using for key-value operations in AsyncMap
type Pair struct {
	First  interface{}
	Second interface{}
}

// AsyncMap provides asynchronous and thread-safe access to map.
// All funcs of this class return channels to interact with map.
type AsyncMap struct {
	cache    map[interface{}]interface{}
	commands chan *asyncCommand
	closed   bool
}

// NewAsyncMap creates new async map.
// cache indicates a size of channel of interim values.
func NewAsyncMap(cache int) *AsyncMap {
	ret := new(AsyncMap)
	ret.cache = make(map[interface{}]interface{})
	ret.commands = make(chan *asyncCommand, cache)
	go func() {
		for cmd := range ret.commands {
			switch cmd.typ {
			case asyncCommandGet:
				if value, ok := ret.cache[cmd.key]; ok {
					cmd.ch <- value
				} else {
					close(cmd.ch)
				}
			case asyncCommandSet:
				ret.cache[cmd.key] = <-cmd.ch
			case asyncCommandDelete:
				delete(ret.cache, cmd.key)
				cmd.ch <- Empty{}
			case asyncCommandLen:
				cmd.ch <- len(ret.cache)
			case asyncCommandList:
				for k, v := range ret.cache {
					cmd.ch <- Pair{k, v}
				}
				close(cmd.ch)
			}
		}
		ret.cache = nil
	}()
	return ret
}

func (am *AsyncMap) checkClosed() {
	if am.closed {
		panic("Async map is closed")
	}
}

// Set returns channel to set the value.
func (am *AsyncMap) Set(key interface{}) chan<- interface{} {
	am.checkClosed()
	ch := make(chan interface{}, 1)
	am.commands <- &asyncCommand{ch, asyncCommandSet, key}
	return ch
}

// Get return channel to get the value.
// If channel closed then map don't have value by the key.
func (am *AsyncMap) Get(key interface{}) <-chan interface{} {
	am.checkClosed()
	ch := make(chan interface{}, 1)
	am.commands <- &asyncCommand{ch, asyncCommandGet, key}
	return ch
}

// Delete ask to delete the element and returns channel which indicates that the element is deleted.
func (am *AsyncMap) Delete(key interface{}) <-chan Empty {
	am.checkClosed()
	ich := make(chan interface{}, 1)
	ch := make(chan Empty, 1)
	go func() {
		ch <- (<-ich).(Empty)
	}()

	am.commands <- &asyncCommand{ich, asyncCommandDelete, key}
	return ch
}

// Len returns channel to get the len of map.
func (am *AsyncMap) Len() <-chan int {
	am.checkClosed()
	ich := make(chan interface{}, 1)
	ch := make(chan int, 1)
	go func() {
		ch <- (<-ich).(int)
	}()
	am.commands <- &asyncCommand{ich, asyncCommandLen, 0}
	return ch
}

// List returns channel for get key-value pairs from map
func (am *AsyncMap) List() <-chan Pair {
	am.checkClosed()
	len := <-am.Len()
	ich := make(chan interface{}, len)
	ch := make(chan Pair, len)
	go func() {
		for p := range ich {
			ch <- p.(Pair)
		}
		close(ch)
	}()
	am.commands <- &asyncCommand{ich, asyncCommandList, 0}
	return ch
}

// Close closes internal channels and finish work goroutine.
func (am *AsyncMap) Close() {
	am.closed = true
	close(am.commands)
}
