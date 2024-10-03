package termites_web

import (
	"encoding/json"
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

	graph.ConnectTo(stateNode.Out, stateTracker.StateIn)
	graph.ConnectTo(connectionsNode.Out, stateTracker.ConnectionIn)

	connectionsNode.Send <- ClientConnection{
		ConnType: 0,
		Id:       "test",
	}
}
