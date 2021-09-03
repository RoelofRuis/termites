package main

import (
	"github.com/RoelofRuis/termites/termites-debug"
)

// Testing the debugger by letting it introspect its own graph..!
func main() {
	graph := debug.NewDebugger(4242).GetGraph()
	debug.Debug(graph, 4243)

	<-graph.Close
}
