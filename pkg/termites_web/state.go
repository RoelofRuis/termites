package termites_web

import (
	"encoding/json"
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
)

// This is a test design with mutable state and diffs that can be sent via websocket.

type StateTracker struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort
	// TODO: out port for state

	serializedState json.RawMessage
}

func NewStateTracker() *StateTracker {
	builder := termites.NewBuilder("StateTracker")

	v := &StateTracker{
		ConnectionIn: termites.NewInPort[ClientConnection](builder, "Connection"),
		StateIn:      termites.NewInPort[json.RawMessage](builder, "State"),

		serializedState: nil,
	}

	builder.OnRun(v.Run)

	return v
}

func (v *StateTracker) Run(c termites.NodeControl) error {
	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			c.LogInfo(fmt.Sprintf("%+v\n", msg))
			// Send full state once, to specific client only!
			// We can know the ClientId, hub should be able to distinguish with special message field.

		case msg := <-v.StateIn.Receive():
			c.LogInfo(fmt.Sprintf("%+v\n", msg))
		}
	}
}
