# Termites

Termites is a flow-based programming framework.

It aims for easy separable and inspectable components, treating each as a separate *node* in a directed *graph*. The
nodes communicate by sending messages to each other over *connections* established between their *ports*.

The `termites` core module provides these basic building blocks.
Optionally, you can use the web components provided in `termites_web` and the debugging tools from `termites_dbg`.

### Development status

The library API is currently still unstable and does not yet serve a v1.

## Examples

See the [examples](termites/examples) folder for example implementations.

## Module `termites`

The termites core module contains an API for creating processing graphs.

This is done by defining `Nodes`, that contain the behaviour of your program. Then a `Graph` is used to tie everything
together and start the processing.

### Usage

When writing a constructor for a node, first create a builder using `termites.NewBuilder`.

Pass this builder when constructing ports for your node, using the `termites.NewInPort` and `termites.NewOutPort`
functions.

```go
package yourpackage

import "github.com/RoelofRuis/termites/pkg/termites"

type YourNode struct {
	In  *termites.InPort  // Receive messages through InPort
	Out *termites.OutPort // Send messages through OutPort
	// ... more ports ...
	// ... other (private) fields ...
}

func NewYourNode() *YourNode {
	builder := termites.NewBuilder("Your Node")

	n := &YourNode{
		In:  termites.NewOutPort[int](builder),
		Out: termites.NewOutPort[int](builder),
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

After you have defined your nodes, connect their ports in a graph:

```go
package main

import "github.com/RoelofRuis/termites/pkg/termites"
import "yourpackage"

func main() {
	graph := termites.NewGraph()

	nodeA := yourpackage.NewYourNode()
	nodeB := yourpackage.NewYourNode()

	graph.Connect(nodeA.Out, nodeB.In)

	graph.Wait()
}
```

### Connection configuration

When connecting nodes, connection options can be passed to configure a single connection.

```golang
var connectionOptions []ConnectionOption
graph.Connect(nodeA.Out, nodeB.In, connectionOptions...) 
```

#### Adding an adapter

The `Via`/`ViaNamed` options insert an adapter function between sender and receiver. 
The adapter is used to transform the sender data so it adheres to the receiver format.

#### Configure the mailbox

To allow for fine-grained control over how messages are handled upon receiving, use the `WithMailbox` connection option.
See the available `MailboxOption` to find out what control is possible.

### Graph configuration

Several configuration options can be passed to the graph constructor:

#### Defer starting the processes until Graph.Wait is called

```golang
graph := termites.NewGraph(termites.DeferredStart())

// Connect the nodes without them directly running

graph.Wait() // Now all nodes/connections will be started together.
```

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

## Module `termites_web`

The web module contains components for easy interaction with the web, mainly through it's websocket component.

### Usage

#### Server side interaction

Create a component that can send and or receive client messages on its ports.

```golang
package yourpackage

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"fmt"
)

type WebProcessor struct {
	In  *termites.InPort
	Out *termites.OutPort
}

func NewWebProcessor() *WebProcessor {
	// The standard termites node setup 
	builder := termites.NewBuilder("WebProcessor")

	n := &WebSender{
		In:  termites.NewInPort[termites_web.ClientMessage](builder),
		Out: termites.NewOutPort[termites_web.ClientMessage](builder),
	}

	builder.OnRun(n.Run)

	return n
}

func (n *YourNode) Run(_ termites.NodeControl) error {
	// In this example, web sender reads a message, sends a message and then immediately quits.
	// Here you could do all kinds of conditional processing, react to other incoming messages, etc.

	msg := <-n.In.Receive()
	fmt.Printf("Received: %+v\n", msg)
	
	n.Out.Send(termites_web.NewClientMessage("your/custom/topic", YourData{42}))

	return nil
}

// YourData can be any data that should be sent to the client.
type YourData struct {
	Count int `json:"count"`
}

```

Configure a connector, bind to a router with your custom router logic attached. Then load and connect the pre-served
javascript code or integrate with your own client side code to send and retrieve data.

```golang
package main

import "github.com/RoelofRuis/termites/pkg/termites"
import "github.com/RoelofRuis/termites/pkg/termites/termites_web"
import "github.com/gorilla/mux"
import "yourpackage"

func main() {
	graph := termites.NewGraph()

	connector := termites_web.NewConnector(graph)

	router := mux.NewRouter()
	router.Path("/ws").Methods("GET").Handler(connector)
	
	// ... other custom (router) setup logic ...

	processor := yourpackage.NewWebProcessor()
	graph.Connect(connector.Hub.OutToApp, processor.In)
	graph.Connect(processor.Out, connector.Hub.InFromApp)

	graph.Wait()
}
```

#### Client side interaction

#### Via the pre-served javascript

In your HTML page include:

```html

<script type="text/javascript" src="/embedded/connect.js"></script>
```

And call:

```html

<script>
    connector.connect();

    connector.subscribe(function (msg) {
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

## Module `termites_dbg`

The debug module contains a web debugger that can be hooked into a graph to inspect it.

### Usage

Initialize it by passing it as an option on graph creation. The debugger will be started on port `4242` by default.

```golang
package main

import "github.com/RoelofRuis/termites/pkg/termites"
import "github.com/RoelofRuis/termites/pkg/termites_dbg"

func main() {
	graph := termites.NewGraph(termites_dbg.WithDebugger())

	// ... graph setup code ...
}
```

### Debugger Configuration

Configuration options can be passed to the debugger Graph Option.

#### Change the debugger port

```golang
termites_dbg.WithDebugger(termites_dbg.OnHttpPort(1234))
```
