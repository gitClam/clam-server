package event

import (
	"clam-server/serverlogger"
	"clam-server/utils/slices"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type dispatcher struct {
	event2HandlerMap map[string][]Handler
	lock             sync.Mutex
}

var p *dispatcher
var once sync.Once

func Dispatcher() *dispatcher {
	once.Do(func() {
		p = &dispatcher{make(map[string][]Handler, 10), sync.Mutex{}}
	})
	return p
}

func (d *dispatcher) AddListener(e Event, eh Handler) {
	name := reflect.TypeOf(e).Name()
	v, _ := d.event2HandlerMap[name]
	// todo 粒度太大（不同eh注册是可以并发的）
	d.lock.Lock()
	defer d.lock.Unlock()
	if slices.ContainsInSlice(v, eh) {
		serverlogger.Warn("Handler registered", zap.String("name", reflect.TypeOf(eh).Name()))
		return
	}
	d.event2HandlerMap[name] = append(d.event2HandlerMap[name], eh)
}
func (d *dispatcher) Dispatch(e Event) {
	name := reflect.TypeOf(e).Name()
	handlers, ok := d.event2HandlerMap[name]
	if !ok {
		return
	}
	for _, handler := range handlers {
		err := handler.Handle(e)
		if err != nil {
			serverlogger.Warn("handle Event err", zap.String("name", name), zap.Error(err))
		}
	}
}
