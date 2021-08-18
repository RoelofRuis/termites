package debug

import (
	"time"

	"github.com/RoelofRuis/termites/termites-core"
)

func WithDebugger(httpPort int) termites.GraphOptions {
	d := ConfigureDebugGraph(termites.NewGraph(termites.NonblockingRun()), httpPort)

	return termites.AddHook(d)
}

func ConfigureDebugGraph(graph *termites.Graph, httpPort int) *debugger {
	// Input for Refs
	nodeRefReceiver := newRefReceiver()

	// Visualizer
	visualizer := NewVisualizer()
	graph.ConnectTo(nodeRefReceiver.RefsOut, visualizer.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// Web UI
	webUI := NewWebController(httpPort)
	graph.ConnectTo(visualizer.PathOut, webUI.PathIn)
	graph.ConnectTo(nodeRefReceiver.RefsOut, webUI.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	return &debugger{
		refReceiver: nodeRefReceiver,
		graph:       graph,
	}
}

type debugger struct {
	refReceiver *refReceiver
	graph       *termites.Graph
}

func (d *debugger) Setup(registry termites.NodeRegistry) {
	d.graph.Run()

	for _, n := range registry.Iterate() {
		n.SetNodeRefChannel(d.refReceiver.refChan)
	}
}

func (d *debugger) Teardown() {
	d.graph.Shutdown()
}
