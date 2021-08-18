package termites

import (
	"testing"
)

func TestCannotRunTwice(t *testing.T) {
	g := NewGraph(NonblockingRun())

	defer func() { recover() }()

	g.Run()
	g.Run()

	t.Errorf("no panic upon running twice") // expected not to be reached because of panic
}

func TestCanRunGraphMultipleTimes(t *testing.T) {
	g := NewGraph(NonblockingRun())

	g.Run()
	g.Shutdown()
	g.Run()
	g.Shutdown()
}

type TestHook struct {
	setupCalls    uint8
	teardownCalls uint8
}

func (h *TestHook) Setup(r NodeRegistry) { h.setupCalls += 1 }
func (h *TestHook) Teardown()            { h.teardownCalls += 1 }

func TestHooks(t *testing.T) {
	testHook := &TestHook{}
	g := NewGraph(AddHook(testHook), NonblockingRun())

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

	g.Run()
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
