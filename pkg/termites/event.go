package termites

import (
	"fmt"
)

type EventType string

type Event struct {
	// Source Identifier TODO: add this later?
	Type EventType
	Data interface{}
}

const (
	InfoLog  EventType = "log/info"
	ErrorLog EventType = "log/error"
	PanicLog EventType = "log/panic"

	NodeRefUpdated EventType = "node/ref-updated"
	NodeStopped    EventType = "node/stopped"

	MessageSent EventType = "message/sent"

	Kill             EventType = "teardown/kill"
	RegisterTeardown EventType = "teardown/register"
	Exit             EventType = "teardown/exit"
)

var InvalidEventError = fmt.Errorf("invalid event")

type InfoLogEvent struct {
	Message string
}

type ErrorLogEvent struct {
	Message string
	Error   error
}

type PanicLogEvent struct {
	Message       string
	Stack         string
	RecoveredData any
}

func LogInfo(msg string) Event {
	return Event{Type: InfoLog, Data: InfoLogEvent{Message: msg}}
}

func LogError(msg string, err error) Event {
	return Event{Type: ErrorLog, Data: ErrorLogEvent{Message: msg, Error: err}}
}

func LogPanic(msg string, stack string, recoveredData any) Event {
	return Event{Type: PanicLog, Data: PanicLogEvent{Message: msg, Stack: stack, RecoveredData: recoveredData}}
}

type RegisterTeardownEvent struct {
	Name string
	F    func(control TeardownControl) error
}

type NodeUpdatedEvent struct {
	Ref NodeRef
}

type NodeStoppedEvent struct {
	Id NodeId
}

type MessageSentEvent struct { // TODO: use ID's
	FromName     string `json:"from_name"`
	FromPortName string `json:"from_port_name"`
	ToName       string `json:"to_name"`
	ToPortName   string `json:"to_port_name"`
	AdapterName  string `json:"adapter_name"`
	Error        string `json:"error"`
}
