# Termites

Termites is a reactive/dataflow framework. 

It aims for easy separable and inspectable components, treating each as a separate *node* in a directed *graph*. The nodes communicate by sending messages to each other over *connections* established between their *ports*.

## Examples

See the examples folder for small example implementations.

## Module `termites`

The termites core module contains an API for creating processing graphs.

This is done by defining `Nodes`, that contain the behaviour of your program. Then a `Graph` is used to tie everything together and start the processing.

### Usage

Create a node using a `builder`

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

Connect your nodes ports in a graph:

```go
package main

import "github.com/RoelofRuis/termites/termites"
import "yourpackage"

func main() {
	graph := termites.NewGraph()

	nodeA := yourpackage.NewYourNode()
	nodeB := yourpackage.NewYourNode()

	graph.ConnectTo(nodeA.Out, nodeB.In)

	graph.Wait()
}
```

### Graph configuration

Several configuration options can be passed to the graph constructor:

#### Attach a console logger
```golang
termites.NewGraph(termites.WithConsoleLogger())
```

#### Close a resource on graph teardown
```golang
closeable io.Closer := YourClosable()
termites.NewGraph(termites.CloseOnShutdown("resource name", closeable))
```

#### Add a custom event subscriber
```golang
sub termites.EventSubscriber := YourSubscriber()
termites.NewGraph(termites.WithEventSubscriber(sub))
```

#### Name the graph
```golang
termites.NewGraph(termites.Named("graph-name"))
```

#### Prevent sigterm handler from being attached
By default, a sigterm handler is attached to the graph that fires an event on sigterm. Passing this option prevents this event from being fired.
```golang
termites.NewGraph(termites.WithoutSigtermHandler())
```

## Module `termites_dbg`

The debug module contains a powerful web debugger that can be hooked into a graph to inspect it.

### Usage

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

## Module `termites_web`

The web module contains components for easy interaction with the web, mainly through it's websocket graph component.

> This module is currently still an unstable work in progress
