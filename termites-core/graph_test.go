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

type TestObserver struct {
	registerCalls uint8
	teardownCalls uint8
}

func (h *TestObserver) Name() string {
	return "Test Observer"
}

func (h *TestObserver) OnNodeRegistered(Node) { h.registerCalls += 1 }
func (h *TestObserver) OnGraphTeardown() { h.teardownCalls += 1 }

func TestObservers(t *testing.T) {
	testHook := &TestObserver{}
	g := NewGraph(AddObserver(testHook))

	if testHook.registerCalls != 0 {
		t.Errorf("expected OnNodeRegistered to be called zero times")
	}
	if testHook.teardownCalls != 0 {
		t.Errorf("expected OnGraphTeardown to be called zero times")
	}

	node := NewInspectableIntNode("A")
	g.ConnectTo(node.Out, node.In)

	if testHook.registerCalls != 1 {
		t.Errorf("expected OnNodeRegistered to be called once")
	}
	if testHook.teardownCalls != 0 {
		t.Errorf("expected OnGraphTeardown to be called zero times")
	}

	g.Shutdown()
	if testHook.registerCalls != 1 {
		t.Errorf("expected setup to be called once")
	}
	if testHook.teardownCalls != 1 {
		t.Errorf("expected OnGraphTeardown to be called once")
	}
}
