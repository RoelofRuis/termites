package termites

import (
	"sync"
)

type EventSubscriber interface {
	SetEventBus(m EventBus)
}

type EventSender interface {
	Send(e Event)
}

type EventBus interface {
	EventSender
	Subscribe(t EventType, f func(Event) error)
}

type eventBus struct {
	subscriptionLock *sync.RWMutex
	subscriptions    map[EventType][]func(event Event) error
	eventChan        chan Event
}

func NewEventBus() *eventBus {
	bus := &eventBus{
		subscriptionLock: &sync.RWMutex{},
		subscriptions:    make(map[EventType][]func(event Event) error),
		eventChan:        make(chan Event, 1000),
	}

	go func() {
		for e := range bus.eventChan {
			bus.subscriptionLock.RLock()
			for _, s := range bus.subscriptions[e.Type] {
				err := s(e)
				if err != nil {
					bus.Send(LogErrorEvent("unable to notify subscriber", err))
				}
			}
			bus.subscriptionLock.RUnlock()
		}
	}()

	return bus
}

func (m *eventBus) Subscribe(t EventType, f func(Event) error) {
	m.subscriptionLock.Lock()
	_, has := m.subscriptions[t]
	if !has {
		m.subscriptions[t] = []func(event Event) error{}
	}
	m.subscriptions[t] = append(m.subscriptions[t], f)
	m.subscriptionLock.Unlock()
}

func (m *eventBus) Send(e Event) {
	// TODO: add timeout? If we can't send here, we are in serious trouble
	m.eventChan <- e
}
