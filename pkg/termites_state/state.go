package termites_state

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

type StateMessage struct {
	Key  string
	Data json.RawMessage
}

// StateStore receives state messages on its StateIn port. These messages can be from different sources, and are
// distinguished by key. The store stores the whole state and creates JSON merge patches from these state change that
// are sent via its PatchOut port.
type StateStore struct {
	StateIn  *termites.InPort
	PatchOut *termites.OutPort

	statesByKey map[string]json.RawMessage
}

func NewStateStore() *StateStore {
	builder := termites.NewBuilder("StateStore")

	s := &StateStore{
		StateIn:  termites.NewInPort[StateMessage](builder, "State"),
		PatchOut: termites.NewOutPort[json.RawMessage](builder, "Patch"),

		statesByKey: make(map[string]json.RawMessage),
	}

	builder.OnRun(s.Run)

	return s
}

func (s *StateStore) Run(c termites.NodeControl) error {
	for msg := range s.StateIn.Receive() {
		stateMessage := msg.Data.(StateMessage)

		original, err := json.Marshal(s.statesByKey)
		if err != nil {
			c.LogError("failed to marshal state", err)
			continue
		}

		s.statesByKey[stateMessage.Key] = stateMessage.Data

		modified, err := json.Marshal(s.statesByKey)
		if err != nil {
			c.LogError("failed to marshal state", err)
			continue
		}

		patch, err := jsonpatch.CreateMergePatch(original, modified)
		if err != nil {
			c.LogError("failed to create merge patch", err)
			continue
		}

		s.PatchOut.Send(json.RawMessage(patch))
	}

	return nil
}
