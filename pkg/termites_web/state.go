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

type State struct {
	ConnectionIn *termites.InPort
	In           *termites.InPort
	MessageOut   *termites.OutPort

	fullState map[string]json.RawMessage
}

func NewState() *State {
	builder := termites.NewBuilder("State")

	t := &State{
		ConnectionIn: termites.NewInPort[ClientConnection](builder),
		In:           termites.NewInPort[StateMessage](builder),
		MessageOut:   termites.NewOutPort[ClientMessage](builder),

		fullState: make(map[string]json.RawMessage),
	}

	builder.OnRun(t.Run)

	return t
}

func (v *State) Run(c termites.NodeControl) error {
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

		case msg := <-v.In.Receive():
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

// MarshalState adapts any data to be wrapped as a JSON encoded state message.
// Use with the termites.Via connection option to set it as state adapter.
func MarshalState(key string) func(in interface{}) (StateMessage, error) {
	return func(in interface{}) (StateMessage, error) {
		data, err := json.Marshal(in)
		if err != nil {
			return StateMessage{}, err
		}
		return StateMessage{Key: key, Data: data}, nil
	}
}
