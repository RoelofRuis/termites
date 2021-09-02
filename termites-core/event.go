package termites

import "fmt"

type EventType string

type Event struct {
	// Source Identifier TODO: add this later?
	Type   EventType
	Data   interface{}
}

const (
	Log            EventType = "log"
	GraphTeardown  EventType = "graph-teardown"
	NodeRegistered EventType = "node-registered"
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
