package debug

import (
	"time"

	"github.com/RoelofRuis/termites/termites-core"
)

func Debug(graph *termites.Graph, httpPort int) {
	debugger := NewDebugger(httpPort)
	graph.Subscribe(debugger)
}

func NewDebugger(httpPort int) *debugger {
	graph := termites.NewGraph(
		termites.Named("Termites Debugger"),
		termites.WithoutSigtermHandler(),
	)

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

func (d *debugger) GetGraph() *termites.Graph {
	return d.graph
}

func (d *debugger) SetEventBus(b *termites.EventBus) {
	b.Subscribe(termites.NodeRefUpdated, d.OnNodeRefUpdated)
	b.Subscribe(termites.GraphTeardown, d.OnGraphTeardown)
}

func (d *debugger) OnNodeRefUpdated(e termites.Event) error {
	n, ok := e.Data.(termites.NodeUpdatedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	d.refReceiver.refChan <- n.Ref
	return nil
}

func (d *debugger) OnGraphTeardown(_ termites.Event) error {
	d.graph.Shutdown()
	return nil
}
