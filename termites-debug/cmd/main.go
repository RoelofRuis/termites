package main

import (
	"github.com/RoelofRuis/termites/termites-core"
	"github.com/RoelofRuis/termites/termites-debug"
)

// Testing the debugger by letting it introspect its own graph..!
func main() {
	graph := termites.NewGraph(debug.WithDebugger(4242)) // TODO: now we need the runner, is this always the case?
	_ = debug.ConfigureDebugGraph(graph, 4243)

	<-graph.Close
}
