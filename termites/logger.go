package termites

import (
	"fmt"
	"log"
)

// TODO: properly format messages

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) SetEventBus(m EventBus) {
	m.Subscribe(Log, l.OnLog)
	m.Subscribe(MessageSent, l.OnMessageSent)
}

func (l *ConsoleLogger) OnLog(e Event) error {
	ev, ok := e.Data.(LogEvent)
	if !ok {
		return InvalidEventError
	}

	log.Printf("LOG: %s", ev.Message)
	return nil
}

func (l *ConsoleLogger) OnMessageSent(e Event) error {
	ev, ok := e.Data.(MessageSentEvent)
	if !ok {
		return InvalidEventError
	}

	connection := formatConnection(ev)

	if ev.Error != nil {
		log.Printf("MESSAGE (ERROR) %s: %s", connection, ev.Error)
	}

	log.Printf("MESSAGE %s", connection)
	return nil
}

func formatConnection(ref MessageSentEvent) string {
	adapterString := ""
	if ref.AdapterName != "" {
		adapterString = fmt.Sprintf("(%s) -> ", ref.AdapterName)
	}
	ownerString := ""
	if ref.ToName != "" {
		ownerString = fmt.Sprintf("%s:%s", ref.ToName, ref.ToPortName)
	}
	return fmt.Sprintf("[%s:%s -> %s%s]",
		ref.FromName,
		ref.FromPortName,
		adapterString,
		ownerString,
	)
}
