package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

type StateTracker struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort
	MessageOut   *termites.OutPort

	fullState json.RawMessage
}

func NewStateTracker() *StateTracker {
	builder := termites.NewBuilder("StateTracker")

	t := &StateTracker{
		ConnectionIn: termites.NewInPort[ClientConnection](builder, "Connection"),
		StateIn:      termites.NewInPort[json.RawMessage](builder, "State"),
		MessageOut:   termites.NewOutPort[ClientMessage](builder, "Message"),

		fullState: nil,
	}

	builder.OnRun(t.Run)

	return t
}

func (v *StateTracker) Run(c termites.NodeControl) error {
	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			connection := msg.Data.(ClientConnection)
			data, _ := MakeMessage("state/full", v.fullState)
			v.MessageOut.Send(ClientMessage{ClientId: connection.Id, Data: data})

		case msg := <-v.StateIn.Receive():
			mergePatch := msg.Data.(json.RawMessage)
			if v.fullState == nil {
				v.fullState = mergePatch
			} else {
				newState, err := jsonpatch.MergePatch(v.fullState, mergePatch)
				if err != nil {
					c.LogError("Failed to apply merge patch", err)
				}
				v.fullState = newState
			}
			data, _ := MakeMessage("state/patch", mergePatch)
			v.MessageOut.Send(ClientMessage{Data: data})
		}
	}
}
