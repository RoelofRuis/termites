package termites

import (
	"sync"
)

type EventSubscriber interface {
	SetEventBus(m *EventBus)
}

type Sender interface {
	Send(e Event)
}

type LoggerSender interface {
	Sender
	LogInfo(msg string)
	LogError(msg string, err error)
}

type EventBus struct {
	subscriptionLock *sync.RWMutex
	subscriptions    map[EventType][]func(event Event) error
	eventChan        chan Event
}

func NewEventBus() *EventBus {
	bus := &EventBus{
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
					bus.LogError("unable to notify subscriber", err)
				}
			}
			bus.subscriptionLock.RUnlock()
		}
	}()

	return bus
}

func (m *EventBus) Subscribe(t EventType, f func(Event) error) {
	m.subscriptionLock.Lock()
	_, has := m.subscriptions[t]
	if !has {
		m.subscriptions[t] = []func(event Event) error{}
	}
	m.subscriptions[t] = append(m.subscriptions[t], f)
	m.subscriptionLock.Unlock()
}

func (m *EventBus) Send(e Event) {
	// TODO: add timeout? If we can't send here, we are in serious trouble
	m.eventChan <- e
}

func (m *EventBus) LogInfo(msg string) {
	m.Send(Event{
		Type: Log,
		Data: LogEvent{
			Level:   1,
			Message: msg,
			Error:   nil,
		},
	})
}

func (m *EventBus) LogError(msg string, err error) {
	m.Send(Event{
		Type: Log,
		Data: LogEvent{
			Level:   3,
			Message: msg,
			Error:   err,
		},
	})
}
