package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
)

type testUser struct {
	name     string
	Nickname string `json:"nickname,omitempty"`
	Balance  int    `json:"balance"`
}

func (t testUser) Mutate(s *testState) error {
	s.Users[t.name] = t
	return nil
}

type testState struct {
	Users map[string]testUser `json:"users"`
}

func (t testState) Read() (json.RawMessage, error) {
	return json.Marshal(t)
}

func TestStateBroadcaster(t *testing.T) {
	graph := termites.NewGraph()

	broadcaster := NewStateBroadcaster(&testState{Users: make(map[string]testUser)})

	mutations := termites.NewInspectableNode[termites.Mutation[*testState]]("Mutations")
	connectionsNode := termites.NewInspectableNode[ClientConnection]("Connections")
	clientMessages := termites.NewInspectableNode[ClientMessage]("Client Messages")

	graph.Connect(mutations.Out, broadcaster.MutationsIn)
	graph.Connect(connectionsNode.Out, broadcaster.ConnectionIn)
	graph.Connect(broadcaster.MessageOut, clientMessages.In)

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "abc123"}

	// Expect the full state (which is nil at this point) to be sent to client abc123 only
	msg := <-clientMessages.Receive
	if msg.ClientId != "abc123" {
		t.Errorf("Expected client id to be 'abc123', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"topic\":\"state/full\",\"payload\":{\"users\":{}}}")

	// Send an update for the state of Alice
	mutations.Send <- testUser{name: "alice", Balance: 42}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"topic\":\"state/patch\",\"payload\":{\"users\":{\"alice\":{\"balance\":42}}}}")

	// Send an update for the state of Bob
	mutations.Send <- testUser{name: "bob", Balance: 10, Nickname: "Disco Bobby"}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"topic\":\"state/patch\",\"payload\":{\"users\":{\"bob\":{\"balance\":10,\"nickname\":\"Disco Bobby\"}}}}")

	// Send an update for the state of Bob
	mutations.Send <- testUser{name: "bob", Balance: 8, Nickname: "Disco Bobby"}

	// Expect a patch state to be sent to all clients
	msg = <-clientMessages.Receive
	if msg.ClientId != "" {
		t.Errorf("Expected client id to be '', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"topic\":\"state/patch\",\"payload\":{\"users\":{\"bob\":{\"balance\":8}}}}")

	// Send a client connect message
	connectionsNode.Send <- ClientConnection{ConnType: ClientConnect, Id: "xyz789"}

	// Expect the full state to be sent to client xyz789 only
	msg = <-clientMessages.Receive
	if msg.ClientId != "xyz789" {
		t.Errorf("Expected client id to be 'abc123', got '%s'", msg.ClientId)
	}

	jsonMustEqual(t, msg.Data, "{\"topic\":\"state/full\",\"payload\":{\"users\":{\"alice\":{\"balance\":42},\"bob\":{\"nickname\":\"Disco Bobby\",\"balance\":8}}}}")
}

func jsonMustEqual(t *testing.T, actual []byte, expected string) {
	if string(actual) != expected {
		t.Errorf("Expected data to be\n'%s'\ngot\n'%s'", expected, actual)
	}
}
