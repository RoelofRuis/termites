# Termites

Termites is a reactive/dataflow framework. It aims for easy separable and inspectable components by sending messages
between them.

## Module `termites`

The core module contains code for creating processing graphs.

This is done by defining `Nodes`, that contain the behaviour of your program. Then a `Graph` is used to tie everything together and start the processing.

### Example

Create a node using the `builder`

```go
package yourpackage

import "github.com/RoelofRuis/termites/termites"

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

import "github.com/RoelofRuis/termites/termites"
import "yourpackage"

func main() {
	graph := termites.NewGraph()

	nodeA := yourpackage.NewYourNode()
	nodeB := yourpackage.NewYourNode()

	graph.ConnectTo(nodeA.Out, nodeB.In)

	<-graph.Close
}
```

### Graph configuration

To customise graph/library behaviour, options can be passed to the graph constructor:

#### Add a console logger
```golang
termites.NewGraph(termites.WithConsoleLogger())
```

#### Close a resource on graph shutdown
```golang
closeable io.Closer := YourClosable()
termites.NewGraph(termites.CloseOnShutdown(closeable))
```

#### Add a custom event subscriber
```golang
sub termites.EventSubscriber := YourSubscriber()
termites.NewGraph(termites.AddEventSubscriber(sub))
```

#### Name the graph
```golang
termites.NewGraph(termites.Named("graph-name"))
```

#### Prevent sigterm handler from being attached
By default, a sigterm handler is attached to the graph that fires an event on sigterm. Passing this option prevents this detection.
```golang
termites.NewGraph(termites.WithoutSigtermHandler())
```

#### Prevent running from being attached
The runner component takes care of starting all the nodes that get connected in the graph. Passing this option prevents the runner from being attached and the nodes from being started.
```golang
termites.NewGraph(termites.WithoutRunner())
```

## Module `termites_dbg`

The debug module contains a powerful web debugger that can be hooked into a graph for inspection (and is itself a graph
as well..!)

Initialize it by passing it as an option on graph creation.

```go
package main

import "github.com/RoelofRuis/termites/termites"
import "github.com/RoelofRuis/termites/termites_dbg"

func main() {
	graph := termites.NewGraph(termites_dbg.WithDebugger(4242))

	// ... graph setup code
}
```

## Module `termites_ws`

The WS module contains a websocket implementation for easy interaction with browser based applications.

> This module is currently still an unstable work in progress
