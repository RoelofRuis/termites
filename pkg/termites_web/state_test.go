package termites_web

import (
	"encoding/json"
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
)

func TestStateTracker(t *testing.T) {
	graph := termites.NewGraph(
		termites.PrintLogsToConsole(),
	)

	stateTracker := NewStateTracker()

	stateNode := termites.NewInspectableNode[json.RawMessage]("State")
	connectionsNode := termites.NewInspectableNode[ClientConnection]("Connections")
	clientMessages := termites.NewInspectableNode[ClientMessage]("Client Messages")

	graph.ConnectTo(stateNode.Out, stateTracker.StateIn)
	graph.ConnectTo(connectionsNode.Out, stateTracker.ConnectionIn)
	graph.ConnectTo(stateTracker.MessageOut, clientMessages.In)

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "test"}

	// Expect the full state (which is nil at this point) to be sent.
	msg := <-clientMessages.Receive
	if msg.ClientId != "test" {
		t.Errorf("Expected client connection to be 'test', got '%s'", msg.ClientId)
	}
	if msg.Data != nil {
		t.Errorf("Expected data to be nil, got %v", msg.Data)
	}

	// Update the state
	stateNode.Send <- []byte("{\"name\": \"testing\", \"count\": 42}")

	// Expect a patch update to be sent to all receivers
	msg = <-clientMessages.Receive
	fmt.Printf("%+v", msg)
}
