package termites

import "fmt"

type EventType string

type Event struct {
	// Source Identifier TODO: add this later?
	Type EventType
	Data interface{}
}

const (
	Log              EventType = "log/log"
	NodeRegistered   EventType = "node/registered"
	NodeRefUpdated   EventType = "node/ref-updated"
	MessageSent      EventType = "message/sent"
	SysExit          EventType = "sys/exit"
	RegisterTeardown EventType = "sigterm/register-teardown"
)

var InvalidEventError = fmt.Errorf("invalid event")

type LogEvent struct {
	Level   uint8
	Message string
	Error   error
}

func LogInfoEvent(msg string) Event {
	return Event{Type: Log, Data: LogEvent{Level: 1, Message: msg, Error: nil}}
}

func LogErrorEvent(msg string, err error) Event {
	return Event{Type: Log, Data: LogEvent{Level: 3, Message: msg, Error: err}}
}

type RegisterTeardownEvent struct {
	f func() error
}

type NodeRegisteredEvent struct {
	node *node
}

type NodeUpdatedEvent struct {
	Ref NodeRef
}

type MessageSentEvent struct { // TODO: use ID's
	FromName     string
	FromPortName string
	ToName       string
	ToPortName   string
	AdapterName  string
	Data         interface{}
	Error        error
}
