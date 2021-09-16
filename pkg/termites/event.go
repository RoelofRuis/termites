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
	Message string
	Stack   string
}

func LogInfo(msg string) Event {
	return Event{Type: InfoLog, Data: InfoLogEvent{Message: msg}}
}

func LogError(msg string, err error) Event {
	return Event{Type: ErrorLog, Data: ErrorLogEvent{Message: msg, Error: err}}
}

func LogPanic(msg string, stack string) Event {
	return Event{Type: PanicLog, Data: PanicLogEvent{Message: msg, Stack: stack}}
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
	FromName     string
	FromPortName string
	ToName       string
	ToPortName   string
	AdapterName  string
	Error        error
}
