package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
)

func TestStateTracker(t *testing.T) {
	graph := termites.NewGraph()

	stateTracker := NewStateTracker()

	stateNode := termites.NewInspectableNode[json.RawMessage]("State")
	connectionsNode := termites.NewInspectableNode[ClientConnection]("Connections")
	clientMessages := termites.NewInspectableNode[ClientMessage]("Client Messages")

	graph.ConnectTo(stateNode.Out, stateTracker.StateIn)
	graph.ConnectTo(connectionsNode.Out, stateTracker.ConnectionIn)
	graph.ConnectTo(stateTracker.MessageOut, clientMessages.In)

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "abc123"}

	// Expect the full state (which is nil at this point) to be sent to client abc123 only
	msg := <-clientMessages.Receive
	if msg.ClientId != "abc123" {
		t.Errorf("Expected client connection to be 'abc123', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/full\",\"payload\":null}")

	// Update the state
	stateNode.Send <- []byte("{\"name\":\"alice\",\"count\":42}")

	// Expect a patch update to be sent to all receivers
	msg = <-clientMessages.Receive

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/patch\",\"payload\":{\"name\":\"alice\",\"count\":42}}")

	// Update the state again.
	stateNode.Send <- []byte("{\"name\": \"bob\"}")

	// Expect a patch update to be sent to all receivers
	msg = <-clientMessages.Receive

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/patch\",\"payload\":{\"name\":\"bob\"}}")

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "789xyz"}

	// Expect the full state to be sent to client 789xyz only
	msg = <-clientMessages.Receive
	if msg.ClientId != "789xyz" {
		t.Errorf("Expected client connection to be '789xyz', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/full\",\"payload\":{\"name\":\"bob\",\"count\":42}}")
}

func jsonMustEqual(t *testing.T, actual []byte, expected string) {
	if string(actual) != expected {
		t.Errorf("Expected data to be '%s', got '%s'", expected, actual)
	}
}
