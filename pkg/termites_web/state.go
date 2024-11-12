package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"strings"
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
	return NewStateWithInitial(make(map[string]json.RawMessage))
}

func NewStateWithInitial(initialState map[string]json.RawMessage) *State {
	builder := termites.NewBuilder("State")

	t := &State{
		ConnectionIn: termites.NewInPort[ClientConnection](builder),
		In:           termites.NewInPort[StateMessage](builder),
		MessageOut:   termites.NewOutPort[ClientMessage](builder),

		fullState: initialState,
	}

	builder.OnRun(t.Run)

	return t
}

func (v *State) Run(c termites.NodeControl) error {
	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			connection := msg.Data.(ClientConnection)

			clientMessage, err := NewClientMessageFor("state/full", connection.Id, v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			v.MessageOut.Send(clientMessage)

		case msg := <-v.In.Receive():
			stateMessage := msg.Data.(StateMessage)

			oldState, err := json.Marshal(v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			if strings.HasPrefix(stateMessage.Key, "[]") { // handle arrays
				key := strings.TrimPrefix(stateMessage.Key, "[]")
				var list []json.RawMessage
				state, has := v.fullState[key]
				if has {
					err := json.Unmarshal(state, &list)
					if err != nil {
						c.LogError("failed to unmarshal state", err)
						continue
					}
				}
				list = append(list, stateMessage.Data)
				data, err := json.Marshal(list)
				if err != nil {
					c.LogError("failed to marshal list", err)
					continue
				}
				v.fullState[key] = data
			} else {
				v.fullState[stateMessage.Key] = stateMessage.Data
			}

			newState, err := json.Marshal(v.fullState)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			patch, err := jsonpatch.CreateMergePatch(oldState, newState)

			clientMessage, _ := NewClientMessage("state/patch", patch)

			v.MessageOut.Send(clientMessage)
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
