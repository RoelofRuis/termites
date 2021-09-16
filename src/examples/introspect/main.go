package main

import (
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
)

// Letting the debugger introspect its own graph 🤯
func main() {
	// Explicitly create separate debugger, so we can bind it on graph creation with WithEventSubscriber.
	debugger := termites_dbg.NewDebugger()
	graph := termites.NewGraph(
		termites.Named("Termites Debugger"),
		termites.PrintLogsToConsole(),
		termites.WithEventSubscriber(debugger),
	)
	// Manually initialize the debugger.
	termites_dbg.Init(graph, debugger)

	// Await termination.
	graph.Wait()
}
