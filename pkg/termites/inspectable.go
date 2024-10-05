package termites

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// InspectableNode should mainly be used as a testing aid.
// Pick a type A for which to create the node, then send and receive messages via the channel.
// Send on the Panic chanel to cause a panic in the run loop.
// Optionally set the Delay value to introduce a delay between receiving a message and passing it on.
type InspectableNode[A any] struct {
	In      *InPort
	Out     *OutPort
	Send    chan A
	Receive chan A
	Panic   chan struct{}
	Delay   time.Duration
}

func NewInspectableNode[A any](name string) *InspectableNode[A] {
	builder := NewBuilder(name)

	n := &InspectableNode[A]{
		In:      NewInPort[A](builder, "int in"),
		Out:     NewOutPort[A](builder, "int out"),
		Send:    make(chan A),
		Receive: make(chan A, 128),
		Panic:   make(chan struct{}),
		Delay:   -1,
	}

	builder.OnRun(n.Run)
	builder.OnShutdown(n.Shutdown)

	return n
}

func (n *InspectableNode[A]) Run(_ NodeControl) error {
	for {
		select {
		case msg := <-n.In.Receive():
			decoded, is := msg.Data.(A)
			if !is {
				log.Printf("decoding message failed: %s\n", msg.Data)
				panic(fmt.Sprintf("decoding message failed: %+v", msg))
			}
			time.Sleep(n.Delay)
			n.Receive <- decoded
			n.Out.Send(decoded)

		case <-n.Panic:
			panic("handler error")

		case v := <-n.Send:
			n.Out.Send(v)
		}
	}
}

func (n *InspectableNode[A]) ReceiveWithin(d time.Duration) (A, error) {
	timer := time.NewTimer(d)
	select {
	case v := <-n.Receive:
		return v, nil
	case <-timer.C:
		var unit A
		return unit, errors.New("did not receive message within timeout")
	}
}

func (n *InspectableNode[A]) Shutdown(_ TeardownControl) error {
	close(n.Send)
	close(n.Receive)
	close(n.Panic)

	return nil
}
