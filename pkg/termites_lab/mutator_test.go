package termites_lab

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"reflect"
	"testing"
	"time"
)

type testState struct {
	Count int
	Name  string
}

type ChangeName struct {
	NewName string
}

func (c ChangeName) Apply(t *testState) {
	t.Name = c.NewName
}

func TestMutator(t *testing.T) {
	graph := termites.NewGraph()

	initState := &testState{
		Count: 42,
		Name:  "Alice",
	}

	mutator := NewMutator(initState)

	actionSender := termites.NewInspectableNode[Action[testState]]("ActionSender")
	stateReceiver := termites.NewInspectableNode[testState]("StateReceiver")

	graph.ConnectTo(actionSender.Out, mutator.ActionIn)
	graph.ConnectTo(mutator.StateOut, stateReceiver.In)

	msg, err := stateReceiver.ReceiveWithin(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(msg, initState) {
		t.Errorf("expected initialState to be returned")
	}

	actionSender.Send <- ChangeName{"Bob"}

	msg, err = stateReceiver.ReceiveWithin(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Count != initState.Count {
		t.Errorf("expected initialState.Count to be %d, got %d", initState.Count, msg.Count)
	}
	if msg.Name != "Bob" {
		t.Errorf("expected initialState.Name to be Bob, got %s", msg.Name)
	}
}
