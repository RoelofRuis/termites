package termites

import (
	"testing"
)

// TODO: test shutdown twice

type TestObserver struct {
	setupCalls    uint8
	teardownCalls uint8
}

func (h *TestObserver) Name() string {
	return "Test Observer"
}

func (h *TestObserver) OnNodeRegistered(Node) {}
func (h *TestObserver) OnGraphTeardown() { h.teardownCalls += 1 }

func TestHooks(t *testing.T) {
	testHook := &TestObserver{}
	g := NewGraph(AddObserver(testHook))

	if testHook.setupCalls != 0 {
		t.Errorf("expected setup to be called zero times")
	}
	if testHook.teardownCalls != 0 {
		t.Errorf("expected teardown to be called zero times")
	}

	g.Shutdown() // Has no effect before the graph is running

	if testHook.setupCalls != 0 {
		t.Errorf("expected setup to be called zero times")
	}
	if testHook.teardownCalls != 0 {
		t.Errorf("expected teardown to be called zero times")
	}

	if testHook.setupCalls != 1 {
		t.Errorf("expected setup to be called once")
	}
	if testHook.teardownCalls != 0 {
		t.Errorf("expected teardown to be zero times")
	}

	g.Shutdown()
	if testHook.setupCalls != 1 {
		t.Errorf("expected setup to be called once")
	}
	if testHook.teardownCalls != 1 {
		t.Errorf("expected teardown to be called once")
	}
}
