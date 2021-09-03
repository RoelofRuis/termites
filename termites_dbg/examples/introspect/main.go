package main

import (
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
)

// Testing the debugger by letting it introspect its own graph 🤯
func main() {
	graph := termites.NewGraph(termites.Named("Termites Debugger"), termites.WithConsoleLogger())
	debugger := termites_dbg.InitGraph(graph, 4242)
	graph.Subscribe(debugger)

	<-graph.Close
}
