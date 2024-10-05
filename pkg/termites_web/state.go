package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
)

// We can have a separate state holder that just holds and updates a state
// Then have that send a diff to whoever wants to accept it.

// Then have the state tracker hold a copy of the state (built by reading the diffs)
// and send that over the web through connection in and diff messages

// This is a test design with mutable state and diffs that can be sent via websocket.

type StateTracker struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort
	MessageOut   *termites.OutPort

	serializedState json.RawMessage
}

func NewStateTracker() *StateTracker {
	builder := termites.NewBuilder("StateTracker")

	v := &StateTracker{
		ConnectionIn: termites.NewInPort[ClientConnection](builder, "Connection"),
		StateIn:      termites.NewInPort[json.RawMessage](builder, "State"),
		MessageOut:   termites.NewOutPort[ClientMessage](builder, "Message"),

		serializedState: nil,
	}

	builder.OnRun(v.Run)

	return v
}

func (v *StateTracker) Run(c termites.NodeControl) error {
	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			connection := msg.Data.(ClientConnection)
			v.MessageOut.Send(ClientMessage{ClientId: connection.Id, Data: v.serializedState})

		case msg := <-v.StateIn.Receive():
			raw := msg.Data.(json.RawMessage)
			// This is a diff always!
			if v.serializedState == nil {
				v.serializedState = raw
			} else {
				// we need to perform the merge patch
			}
			v.MessageOut.Send(ClientMessage{Data: raw})
		}
	}
}
