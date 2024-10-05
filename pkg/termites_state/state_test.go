package termites_state

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
)

func TestStateStore(t *testing.T) {
	graph := termites.NewGraph()

	stateStore := NewStateStore()

	stateSender := termites.NewInspectableNode[StateMessage]("StateSender")
	patchReceiver := termites.NewInspectableNode[json.RawMessage]("PatchReceiver")

	graph.ConnectTo(stateSender.Out, stateStore.StateIn)
	graph.ConnectTo(stateStore.PatchOut, patchReceiver.In)

	stateSender.Send <- StateMessage{
		Key:  "alice",
		Data: []byte("{\"balance\": 100}"),
	}

	msg := <-patchReceiver.Receive
	if string(msg) != "{\"alice\":{\"balance\":100}}" {
		t.Errorf("Invalid result message")
	}

	stateSender.Send <- StateMessage{
		Key:  "bob",
		Data: []byte("{\"balance\": 60}"),
	}

	msg = <-patchReceiver.Receive
	if string(msg) != "{\"bob\":{\"balance\":60}}" {
		t.Errorf("Invalid result message")
	}

	stateSender.Send <- StateMessage{
		Key:  "alice",
		Data: []byte("{\"balance\": 100, \"name\": \"Alice Allison\"}"),
	}

	msg = <-patchReceiver.Receive
	if string(msg) != "{\"alice\":{\"name\":\"Alice Allison\"}}" {
		t.Errorf("Invalid result message")
	}
}
