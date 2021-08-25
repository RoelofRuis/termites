## Termites

Termites is a reactive/dataflow framework. It aims for easy separable and inspectable components by sending messages
between them.

## Modules

### `termites-core`

The core module contains code for creating processing graphs.

Create a component using the `builder`

```go
package yourpackage

import "github.com/RoelofRuis/termites/termites-core"

type YourNode struct {
	In  *termites.InPort  // Receive messages through InPort
	Out *termites.OutPort // Send messages through OutPort
	// ... more ports
	// ... other (private) fields
}

func NewYourNode() *YourNode {
	builder := termites.NewBuilder("Your Node")

	n := &YourNode{
		In:  builder.InPort("In", 0),
		Out: builder.OutPort("Out", 0),
	}

	builder.OnRun(n.Run)

	return n
}

func (n *YourNode) Run(_ termites.NodeControl) error {
	// any logic can go here, but this kind of construct is useful for streaming processing
	for {
		select {
		case msg := <-n.In.Receive():
			// process message
			n.Out.Send(msg)
		}
	}
}
```

Tie your components together in a graph:

```go
package main

import "github.com/RoelofRuis/termites/termites-core"
import "yourpackage"

func main() {
	graph := termites.NewGraph()

	nodeA := yourpackage.NewYourNode()
	nodeB := yourpackage.NewYourNode()

	graph.ConnectTo(nodeA.Out, nodeB.In)

	graph.Run()
}
```

### `termites-debug`

The debug module contains a powerful web debugger that can be hooked into a graph for inspection (and is itself a graph
as well..!)

Initialize it by passing it to the graph initializer

```go
package main

import "github.com/RoelofRuis/termites/termites-core"
import "github.com/RoelofRuis/termites/termites-debug"

func main() {
	graph := termites.NewGraph(debug.WithDebugger(4242))

	// ... graph setup code

	graph.Run()
}
```

### `termites-ws`

The WS module contains a websocket implementation for easy interaction with browser based applications.

> This module is currently still an unstable work in progress
