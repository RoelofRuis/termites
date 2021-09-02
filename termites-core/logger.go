package termites

import (
	"fmt"
	"log"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) SetEventBus(m *EventBus) {
	m.Subscribe(Log, l.OnMessageLogged)
}

func (l *ConsoleLogger) OnMessageLogged(e Event) error {
	ev, ok := e.Data.(LogEvent)
	if !ok {
		return fmt.Errorf("logger received event [%+v] of invalid type", e)
	}

	log.Printf(ev.Message)
	return nil
}

// TODO: where should this go?
func formatRoute(ref MessageRef) string {
	adapterString := ""
	if ref.adapterName != "" {
		adapterString = fmt.Sprintf("(%s) -> ", ref.adapterName)
	}
	ownerString := ""
	if ref.toName != "" {
		ownerString = fmt.Sprintf("%s:%s", ref.toName, ref.toPortName)
	}
	return fmt.Sprintf("[%s:%s -> %s%s]",
		ref.fromName,
		ref.fromPortName,
		adapterString,
		ownerString,
	)
}
