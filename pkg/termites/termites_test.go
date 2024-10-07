package termites

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConnections(t *testing.T) {
	graph := NewGraph()

	nodeA := NewInspectableNode[int]("Component A")
	nodeB := NewInspectableNode[int]("Component B")
	nodeC := NewInspectableNode[int]("Component C")
	nodeD := NewInspectableNode[int]("Component D")

	graph.Connect(nodeA.Out, nodeB.In)
	graph.Connect(nodeB.Out, nodeC.In)
	graph.Connect(nodeB.Out, nodeD.In)

	nodeA.Send <- 42

	vC := <-nodeC.Receive
	vD := <-nodeD.Receive
	if vC == 42 && vD == 42 {
		t.Log("Both nodes received correct message")
	} else {
		t.Errorf("Incorrect message")
	}
}

func TestDynamicConnections(t *testing.T) {
	graph := NewGraph()

	nodeA := NewInspectableNode[int]("Component A")
	nodeB1 := NewInspectableNode[int]("Component B1")
	nodeB2 := NewInspectableNode[int]("Component B2")

	connB1 := graph.Connect(nodeA.Out, nodeB1.In)

	nodeA.Send <- 42

	vB1 := <-nodeB1.Receive
	if vB1 == 42 {
		t.Log("Node B1 received correct message")
	} else {
		t.Error("Incorrect message")
	}

	graph.Connect(nodeA.Out, nodeB2.In)

	nodeA.Send <- 43

	vB1 = <-nodeB1.Receive
	vB2 := <-nodeB2.Receive
	if vB1 == 43 && vB2 == 43 {
		t.Log("Nodes received correct message")
	} else {
		t.Error("Incorrect message")
	}

	connB1.Disconnect()

	nodeA.Send <- 44

	vB2 = <-nodeB2.Receive
	if vB2 == 44 {
		t.Log("Node B2 received correct message")
	} else {
		t.Error("Incorrect message")
	}
}

func TestAdapter(t *testing.T) {
	graph := NewGraph()

	stringToInt := func(in string) (int, error) {
		if in == "skip" {
			return 0, SkipElement
		}
		i, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			return 0, err
		}
		return int(i), nil
	}

	nodeA := NewInspectableNode[string]("Component A")
	nodeB := NewInspectableNode[int]("Component B")

	graph.Connect(nodeA.Out, nodeB.In, Via(stringToInt))

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

	nodeA := NewInspectableNode[int]("Component A")
	nodeB := NewInspectableNode[int]("Component B")

	graph.Connect(nodeA.Out, nodeB.In)

	nodeA.Panic <- struct{}{}

	time.Sleep(100 * time.Millisecond)
	t.Log("Program has not crashed")
}

func TestTimeout(t *testing.T) {
	graph := NewGraph()

	nodeA := NewInspectableNode[int]("Component A")
	nodeB := NewInspectableNode[int]("Component B")
	nodeB.Delay = 2 * time.Second // Should time out

	graph.Connect(nodeA.Out, nodeB.In)

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

	nodeA := NewInspectableNode[int]("Component A")
	nodeB := NewInspectableNode[int]("Component B")
	nodeB.Delay = 500 * time.Millisecond
	nodeC := NewInspectableNode[int]("Component C")
	nodeC.Delay = 500 * time.Millisecond
	nodeD := NewInspectableNode[int]("Component D")
	nodeD.Delay = 500 * time.Millisecond

	graph.Connect(nodeA.Out, nodeB.In)
	graph.Connect(nodeA.Out, nodeC.In)
	graph.Connect(nodeA.Out, nodeD.In)

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
