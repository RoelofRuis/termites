package main

import (
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
)

// Testing the debugger by letting it introspect its own graph 🤯
func main() {
	debugger := termites_dbg.NewDebugger(4242)
	graph := debugger.GetGraph()
	graph.Subscribe(debugger)
	graph.Subscribe(termites.NewSigtermHandler()) // Debugger uses sigterm from main graph by default

	<-graph.Close
}
