package termites

import "fmt"

type EventType string

type Event struct {
	// Source Identifier TODO: add this later?
	Type EventType
	Data interface{}
}

const (
	Log            EventType = "log"
	GraphTeardown  EventType = "graph/teardown"
	NodeRegistered EventType = "node/registered"
	NodeRefUpdated EventType = "node/ref-updated"
	MessageSent    EventType = "message/sent"
)

var InvalidEventError = fmt.Errorf("invalid event")

type LogEvent struct {
	Level   uint8
	Message string
	Error   error
}

type NodeRegisteredEvent struct {
	Node Node
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
