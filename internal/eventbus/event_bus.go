package eventbus

import (
	"sync"
)

var bus = &eventBus{channels: &sync.Map{}, lock: &sync.Mutex{}}

// Publish Post an event.
func Publish(topic string, args ...any) {
	bus.Publish(topic, args...)
}

// Subscribe to events
func Subscribe(topic string, fn func(args ...any)) {
	bus.Subscribe(topic, fn)
}

type eventBus struct {
	channels *sync.Map
	lock     *sync.Mutex
}

func (s *eventBus) Publish(topic string, args ...any) {
	subs, ok := s.channels.Load(topic)
	if ok {
		subs.(*sync.Map).Range(func(key, value any) bool {
			key.(*handler).handle(args...)
			return true
		})
	}
}

func (s *eventBus) Subscribe(topic string, fn func(args ...any)) {
	subs, exists := s.channels.Load(topic)
	if !exists {
		bus.lock.Lock()
		defer bus.lock.Unlock()
		subs, exists = s.channels.Load(topic)
		if !exists {
			subs = &sync.Map{}
			s.channels.Store(topic, subs)
		}
	}
	subs.(*sync.Map).Store(&handler{fn}, true)
}

type handler struct {
	handle func(args ...any)
}
