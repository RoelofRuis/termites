package termites

import (
	"testing"
)

type testCloser struct {
	done chan interface{}
}

func (c *testCloser) Close() error {
	close(c.done)
	return nil
}

func TestKill(t *testing.T) {
	closer := &testCloser{done: make(chan interface{})}
	graph := NewGraph(CloseOnTeardown("test", closer))

	graph.Close()

	<-closer.done
	t.Log("Closer was closed correctly")
}

type TestSubscriber struct {
	registerCalls uint8
	teardownCalls uint8
}

func (s *TestSubscriber) SetEventBus(e EventBus) {

}

func (s *TestSubscriber) OnNodeRegistered(_ Event) error {
	s.registerCalls += 1
	return nil
}

func TestObservers(t *testing.T) {
	// TODO: this test must be rewritten using timeouts
}
