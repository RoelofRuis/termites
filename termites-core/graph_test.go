package termites

import (
	"testing"
)

func TestShutdownTwice(t *testing.T) {
	g := NewGraph()

	g.Shutdown()
	g.Shutdown()

	<-g.Close
}

type TestSubscriber struct {
	registerCalls uint8
	teardownCalls uint8
}

func (t *TestSubscriber) SetEventBus(e *EventBus) {
	e.Subscribe(GraphTeardown, t.OnGraphTeardown)
	e.Subscribe(NodeRegistered, t.OnNodeRegistered)
}

func (h *TestSubscriber) OnNodeRegistered(_ Event) error {
	h.registerCalls += 1
	return nil
}

func (h *TestSubscriber) OnGraphTeardown(_ Event) error {
	h.teardownCalls += 1
	return nil
}

func TestObservers(t *testing.T) {
	testSubscriber := &TestSubscriber{}
	g := NewGraph(AddEventSubscriber(testSubscriber))

	if testSubscriber.registerCalls != 0 {
		t.Errorf("expected OnNodeRegistered to be called zero times")
	}
	if testSubscriber.teardownCalls != 0 {
		t.Errorf("expected OnGraphTeardown to be called zero times")
	}

	node := NewInspectableIntNode("A")
	g.ConnectTo(node.Out, node.In)

	if testSubscriber.registerCalls != 1 {
		t.Errorf("expected OnNodeRegistered to be called once")
	}
	if testSubscriber.teardownCalls != 0 {
		t.Errorf("expected OnGraphTeardown to be called zero times")
	}

	g.Shutdown()
	if testSubscriber.registerCalls != 1 {
		t.Errorf("expected setup to be called once")
	}
	if testSubscriber.teardownCalls != 1 {
		t.Errorf("expected OnGraphTeardown to be called once")
	}
}
