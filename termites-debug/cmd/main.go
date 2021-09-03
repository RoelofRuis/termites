package main

import (
	"github.com/RoelofRuis/termites/termites-core"
	"github.com/RoelofRuis/termites/termites-debug"
)

// Testing the debugger by letting it introspect its own graph 🤯
func main() {
	debugger := debug.NewDebugger(4242)
	graph := debugger.GetGraph()
	graph.Subscribe(debugger)
	graph.Subscribe(termites.NewSigtermHandler()) // Debugger uses sigterm from main graph by default

	<-graph.Close
}
