package termites

import (
	"fmt"
	"log"
)

type ConsoleLogger struct {
	printLogs     bool
	printMessages bool
}

func NewConsoleLogger(printLogs bool, printMessages bool) *ConsoleLogger {
	return &ConsoleLogger{
		printLogs:     printLogs,
		printMessages: printMessages,
	}
}

func (l *ConsoleLogger) SetEventBus(m EventBus) {
	if l.printLogs {
		m.Subscribe(InfoLog, l.OnLogInfo)
		m.Subscribe(ErrorLog, l.OnLogError)
		m.Subscribe(PanicLog, l.OnLogPanic)
	}

	if l.printMessages {
		m.Subscribe(MessageSent, l.OnMessageSent)
	}
}

func (l *ConsoleLogger) OnLogInfo(e Event) error {
	ev, ok := e.Data.(InfoLogEvent)
	if !ok {
		return InvalidEventError
	}

	log.Printf("LOG: %s", ev.Message)
	return nil
}

func (l *ConsoleLogger) OnLogError(e Event) error {
	ev, ok := e.Data.(ErrorLogEvent)
	if !ok {
		return InvalidEventError
	}

	log.Printf("ERROR: %s: %s", ev.Message, ev.Error)
	return nil
}

func (l *ConsoleLogger) OnLogPanic(e Event) error {
	ev, ok := e.Data.(PanicLogEvent)
	if !ok {
		return InvalidEventError
	}

	log.Printf("PANIC: %s", ev.Message)
	log.Printf("----- stack trace -----\n%s\n------ end trace -------", ev.Stack)
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
		return nil
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
