package termites

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
	subscriptions map[EventType][]func(event Event) error
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscriptions: make(map[EventType][]func(event Event) error),
	}
}

func (m *EventBus) Subscribe(t EventType, f func(Event) error) {
	subs, has := m.subscriptions[t]
	if !has {
		subs = []func(event Event) error{}
	}
	subs = append(subs, f)
	m.subscriptions[t] = subs
}

func (m *EventBus) Send(e Event) {
	// TODO: should this be async? Should there be another decoupling?
	for _, s := range m.subscriptions[e.Type] {
		err := s(e)
		if err != nil {
			m.LogError("unable to notify subscriber", err)
		}
	}
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
