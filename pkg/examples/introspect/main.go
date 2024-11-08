package main

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_dbg"
	"github.com/RoelofRuis/termites/pkg/termites_web"
)

// Letting the debugger introspect its own graph 🤯
func main() {
	// Explicitly create separate debugger, so we can bind it on graph creation with WithEventSubscriber.
	debugger := termites_dbg.NewDebugger(
		termites_dbg.OnHttpPort(4242),
		termites_dbg.WithoutMessageTracking(), // Disable message tracking as this would create infinite message loops.
	)

	// Create a new graph
	graph := termites.NewGraph(
		termites.Named("Termites Debugger"),
		termites.PrintLogsToConsole(),
		termites.WithEventSubscriber(debugger),
	)
	// Manually initialize the debugger.
	termites_dbg.Init(graph, debugger)

	// Open a browser http://localhost:4242 in your browser to see the debugger inspecting itself!
	_ = termites_web.RunBrowser("http://localhost:4242")

	// Await termination.
	graph.Wait()
}
