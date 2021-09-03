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

func (t *TestSubscriber) SetEventBus(e EventBus) {
	e.Subscribe(NodeRegistered, t.OnNodeRegistered)
}

func (h *TestSubscriber) OnNodeRegistered(_ Event) error {
	h.registerCalls += 1
	return nil
}

func TestObservers(t *testing.T) {
	// TODO: this test must be rewritten using timeouts
}
