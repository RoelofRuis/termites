package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

type StateMessage struct {
	Key  string
	Data json.RawMessage
}

type StateTracker struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort
	MessageOut   *termites.OutPort

	fullState map[string]json.RawMessage
}

func NewStateTracker() *StateTracker {
	builder := termites.NewBuilder("StateTracker")

	t := &StateTracker{
		ConnectionIn: termites.NewInPort[ClientConnection](builder, "Connection"),
		StateIn:      termites.NewInPort[StateMessage](builder, "State"),
		MessageOut:   termites.NewOutPort[ClientMessage](builder, "Message"),

		fullState: make(map[string]json.RawMessage),
	}

	builder.OnRun(t.Run)

	return t
}

func (v *StateTracker) Run(c termites.NodeControl) error {
	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			connection := msg.Data.(ClientConnection)

			newState, err := json.Marshal(v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			data, _ := WebUpdate("state/full", newState)
			v.MessageOut.Send(ClientMessage{ClientId: connection.Id, Data: data})

		case msg := <-v.StateIn.Receive():
			stateMessage := msg.Data.(StateMessage)

			oldState, err := json.Marshal(v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			v.fullState[stateMessage.Key] = stateMessage.Data

			newState, err := json.Marshal(v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			patch, err := jsonpatch.CreateMergePatch(oldState, newState)

			data, _ := WebUpdate("state/patch", patch)

			v.MessageOut.Send(ClientMessage{Data: data})
		}
	}
}
