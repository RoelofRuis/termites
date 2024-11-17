package termites

import "testing"

const (
	EventA EventType = "Event A"
	EventB EventType = "Event B"
)

type Counter struct {
	counts map[EventType]int
	done   chan interface{}
}

func (c *Counter) Count(e Event) error {
	_, has := c.counts[e.Type]
	if !has {
		c.counts[e.Type] = 0
	}
	c.counts[e.Type]++
	done, ok := e.Data.(bool)
	if ok && done {
		close(c.done)
	}
	return nil
}

func TestEventBus(t *testing.T) {
	bus := newEventBus()

	counter := &Counter{
		counts: make(map[EventType]int),
		done:   make(chan interface{}),
	}

	bus.Subscribe(EventA, counter.Count)
	bus.Subscribe(EventB, counter.Count)

	bus.Send(Event{Type: EventA})
	bus.Send(Event{Type: EventA})
	bus.Send(Event{Type: EventB})
	bus.Send(Event{Type: EventA, Data: true})

	<-counter.done

	if counter.counts[EventA] != 3 {
		t.Errorf("expected to receive EventA %d times, got %d", 3, counter.counts[EventA])
	}
	if counter.counts[EventB] != 1 {
		t.Errorf("expected to receive EventB %d times, got %d", 1, counter.counts[EventB])
	}
}
