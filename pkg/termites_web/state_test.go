package termites_web

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
)

func TestState(t *testing.T) {
	graph := termites.NewGraph()

	state := NewState()

	stateNode := termites.NewInspectableNode[StateMessage]("State")
	connectionsNode := termites.NewInspectableNode[ClientConnection]("Connections")
	clientMessages := termites.NewInspectableNode[ClientMessage]("Client Messages")

	graph.Connect(stateNode.Out, state.In)
	graph.Connect(connectionsNode.Out, state.ConnectionIn)
	graph.Connect(state.MessageOut, clientMessages.In)

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "abc123"}

	// Expect the full state (which is nil at this point) to be sent to client abc123 only
	msg := <-clientMessages.Receive
	if msg.ClientId != "abc123" {
		t.Errorf("Expected client id to be 'abc123', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/full\",\"payload\":{}}")

	// Send an update for the state of Alice
	stateNode.Send <- StateMessage{Key: "alice", Data: []byte("{\"balance\": 42}")}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/patch\",\"payload\":{\"alice\":{\"balance\":42}}}")

	// Send an update for the state of Bob
	stateNode.Send <- StateMessage{Key: "bob", Data: []byte("{\"balance\": 10, \"name\": \"Bobby\"}")}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/patch\",\"payload\":{\"bob\":{\"balance\":10,\"name\":\"Bobby\"}}}")

	// Send an update for the state of Bob
	stateNode.Send <- StateMessage{Key: "bob", Data: []byte("{\"balance\": 8, \"name\": \"Bobby\"}")}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/patch\",\"payload\":{\"bob\":{\"balance\":8}}}")

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "xyz789"}

	// Expect the full state to be sent to client xyz789 only
	msg = <-clientMessages.Receive
	if msg.ClientId != "xyz789" {
		t.Errorf("Expected client id to be 'abc123', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"msg_type\":\"update\",\"content_type\":\"state/full\",\"payload\":{\"alice\":{\"balance\":42},\"bob\":{\"balance\":8,\"name\":\"Bobby\"}}}")
}

func jsonMustEqual(t *testing.T, actual []byte, expected string) {
	if string(actual) != expected {
		t.Errorf("Expected data to be '%s', got '%s'", expected, actual)
	}
}
