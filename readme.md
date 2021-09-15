# Termites

Termites is a reactive/dataflow framework.

It aims for easy separable and inspectable components, treating each as a separate *node* in a directed *graph*. The
nodes communicate by sending messages to each other over *connections* established between their *ports*.

The `termites` core module provides these basic building blocks. Optionally, you can use the debugging tools from `termites_dbg` and/or the web components provided in `termites_web`.

## Examples

See the [examples](examples) folder for small example implementations.

## Module `termites`

The termites core module contains an API for creating processing graphs.

This is done by defining `Nodes`, that contain the behaviour of your program. Then a `Graph` is used to tie everything
together and start the processing.

### Usage

Create a node using a `builder`

```go
package yourpackage

import "github.com/RoelofRuis/termites/termites"

type YourNode struct {
	In  *termites.InPort  // Receive messages through InPort
	Out *termites.OutPort // Send messages through OutPort
	// ... more ports ...
	// ... other (private) fields ...
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
	// Any logic can go here, but this kind of construct is useful for streaming processing
	for {
		select {
		case msg := <-n.In.Receive():
			// ... Process the message ...
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

#### Printing logs

```golang
// To show the logs
termites.NewGraph(termites.PrintLogsToConsole())

// To show the messages
termites.NewGraph(termites.PrintMessagesToConsole())
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

By default, a sigterm handler is attached to the graph that fires an event on sigterm. Passing this option prevents this
event from being fired.

```golang
termites.NewGraph(termites.WithoutSigtermHandler())
```

## Module `termites_dbg`

The debug module contains a web debugger that can be hooked into a graph to inspect it.

### Usage

Initialize it by passing it as an option on graph creation.

```golang
package main

import "github.com/RoelofRuis/termites/termites"
import "github.com/RoelofRuis/termites/termites_dbg"

func main() {
	graph := termites.NewGraph(termites_dbg.WithDebugger(4242))

	// ... graph setup code ...
}
```

## Module `termites_web`

The web module contains components for easy interaction with the web, mainly through it's websocket graph component.

### Usage

Configure a connector, bind to a router with your custom router logic attached. Then load and connect the pre-served
javascript code to send and retrieve data.

```golang
package main

import "github.com/RoelofRuis/termites/termites"
import "github.com/RoelofRuis/termites/termites/termites_web"
import "github.com/gorilla/mux"

func main() {
	graph := termites.NewGraph()

	connector := termites_web.NewConnector(graph)

	router := mux.NewRouter()
	connector.Bind(router)

	// ... your custom (router) setup logic ...

	graph.Wait()
}
```

In your HTML page include:

```html

<script type="text/javascript" src="/embedded/connect.js"></script>
```

And call:

```html

<script>
    connector.connect();

    connector.subscribe(function (tpe, data) {
        // ... connect to your front-end logic ...
    })
</script>
```

### Browser control

Start a browser window opening the given url using a function or a `termites` component.

```golang
// Function syntax
termites_web.RunBrowser("localhost:8000")

// Component syntax
manager := termites_web.NewBrowserManager("localhost:8000")
```
