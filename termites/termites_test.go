package termites

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

type InspectableIntNode struct {
	In      *InPort
	Out     *OutPort
	Send    chan int
	Receive chan int
}

func NewInspectableIntNode(name string) *InspectableIntNode {
	builder := NewBuilder(name)

	n := &InspectableIntNode{
		In:      builder.InPort("int in", 0),
		Out:     builder.OutPort("int out", 0),
		Send:    make(chan int),
		Receive: make(chan int, 128),
	}

	builder.OnRun(n.Run)

	return n
}

func (c *InspectableIntNode) Run(_ NodeControl) error {
	for {
		select {
		case msg := <-c.In.Receive():
			decoded := msg.Data.(int)
			c.Receive <- decoded
			c.Out.Send(decoded)

		case v := <-c.Send:
			if v == -1 {
				panic("handler error")
			}
			c.Out.Send(v)
		}
	}
}

type InspectableStringNode struct {
	In  *InPort
	Out *OutPort
}

func NewInspectableStringNode(name string) *InspectableStringNode {
	builder := NewBuilder(name)

	n := &InspectableStringNode{
		In:  builder.InPort("string in", ""),
		Out: builder.OutPort("string out", ""),
	}

	builder.OnRun(n.Run)

	return n
}

func (c *InspectableStringNode) Run(_ NodeControl) error {
	for {
		select {
		case msg := <-c.In.Receive():
			decoded := msg.Data.(string)
			c.Out.Send(decoded)
		}
	}
}

type DelayedIntNode struct {
	In      *InPort
	Out     *OutPort
	Receive chan int
	Delay   time.Duration
}

func NewDelayedIntNode(name string, delay time.Duration) *DelayedIntNode {
	builder := NewBuilder(name)

	n := &DelayedIntNode{
		In:      builder.InPort("int in", 0),
		Out:     builder.OutPort("int out", 0),
		Receive: make(chan int),
		Delay:   delay,
	}

	builder.OnRun(n.Run)

	return n
}

func (c *DelayedIntNode) Run(_ NodeControl) error {
	for {
		select {
		case msg := <-c.In.Receive():
			decoded := msg.Data.(int)
			time.Sleep(c.Delay) // simulate work
			c.Receive <- decoded
			c.Out.Send(decoded)
		}
	}
}

// TESTS

func TestConnections(t *testing.T) {
	graph := NewGraph()

	nodeA := NewInspectableIntNode("Component A")
	nodeB := NewInspectableIntNode("Component B")
	nodeC := NewInspectableIntNode("Component C")
	nodeD := NewInspectableIntNode("Component D")

	graph.ConnectTo(nodeA.Out, nodeB.In)
	graph.ConnectTo(nodeB.Out, nodeC.In)
	graph.ConnectTo(nodeB.Out, nodeD.In)

	nodeA.Send <- 42

	vC := <-nodeC.Receive
	vD := <-nodeD.Receive
	if vC == 42 && vD == 42 {
		t.Log("Both nodes received correct message")
	} else {
		t.Errorf("Incorrect message")
	}
}

func TestAdapter(t *testing.T) {
	graph := NewGraph()

	adapter := NewAdapter(
		"string to int",
		"",
		0,
		func(in interface{}) (interface{}, error) {
			s := in.(string)
			if s == "skip" {
				return nil, nil
			}
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			return int(i), nil
		},
	)

	nodeA := NewInspectableStringNode("Component A")
	nodeB := NewInspectableIntNode("Component B")

	graph.ConnectTo(nodeA.Out, nodeB.In, Via(adapter))

	nodeA.Out.Send("skip")
	nodeA.Out.Send("42")

	res := <-nodeB.Receive
	if res == 42 {
		t.Log("component received transformed result")
	} else {
		t.Errorf("Incorrect message")
	}
}

func TestNodePanic(t *testing.T) {
	graph := NewGraph()

	nodeA := NewInspectableIntNode("Component A")
	nodeB := NewInspectableIntNode("Component B")

	graph.ConnectTo(nodeA.Out, nodeB.In)

	nodeA.Send <- -1

	time.Sleep(100 * time.Millisecond)
	t.Log("Program has not crashed")
}

func TestTimeout(t *testing.T) {
	graph := NewGraph()

	nodeA := NewDelayedIntNode("Component A", 0*time.Second)
	nodeB := NewDelayedIntNode("Component B", 2*time.Second) // Should time out

	graph.ConnectTo(nodeA.Out, nodeB.In)

	go func() {
		nodeA.Out.Send(41)
		nodeA.Out.Send(42)
	}()

	count := 0
	done := false
	timeout := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-nodeB.Receive:
			count++
		case <-timeout.C:
			done = true
		}
		if count == 2 {
			t.Errorf("Received 2 messages, but expected 1 to time out")
			break
		}
		if done {
			t.Log("Received 1 only 1 message, other timed out")
			break
		}
	}
}

func TestAsyncSendTiming(t *testing.T) {
	graph := NewGraph()

	nodeA := NewDelayedIntNode("Component A", 0*time.Second)
	nodeB := NewDelayedIntNode("Component B", 500*time.Millisecond)
	nodeC := NewDelayedIntNode("Component C", 500*time.Millisecond)
	nodeD := NewDelayedIntNode("Component D", 500*time.Millisecond)

	graph.ConnectTo(nodeA.Out, nodeB.In)
	graph.ConnectTo(nodeA.Out, nodeC.In)
	graph.ConnectTo(nodeA.Out, nodeD.In)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		<-nodeB.Receive
		<-nodeB.Receive
		wg.Done()
	}()
	go func() {
		<-nodeC.Receive
		<-nodeC.Receive
		wg.Done()
	}()
	go func() {
		<-nodeD.Receive
		<-nodeD.Receive
		wg.Done()
	}()
	done := make(chan interface{})
	go func() {
		wg.Wait()
		done <- nil
	}()

	nodeA.Out.Send(41)
	nodeA.Out.Send(42)

	maxDur := time.NewTicker(1200 * time.Millisecond)
	select {
	case <-maxDur.C:
		t.Errorf("Test timed out: not all messages received")
		t.FailNow()
	case <-done:
		t.Log("All messages received in time")
	}
}
