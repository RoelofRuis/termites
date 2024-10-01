package termites_dbg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"github.com/gorilla/mux"
)

// WithDebugger initializes a new debugger (represented by its own termites.Graph!) and returns it as an option so it
// can be attached the main program graph.
func WithDebugger(opts ...DebuggerOption) termites.GraphOption {
	dbg := NewDebugger(opts...)

	graph := termites.NewGraph(
		termites.Named("Termites Debugger"),
		termites.WithoutSigtermHandler(),
		termites.WithEventSubscriber(dbg.tempDir),
	)

	Init(graph, dbg)

	return termites.WithEventSubscriber(dbg)
}

// Init initializes the given graph with the given debugger. Prefer to use the termites.GraphOption WithDebugger
// function to initiate a connection.
func Init(graph *termites.Graph, debugger *Debugger) {
	connector := termites_web.NewConnector(graph)

	router := mux.NewRouter()
	connector.Bind(router)

	controller := NewWebController()
	controller.editor = debugger.editor

	router.HandleFunc("/", controller.HandleIndex)
	router.HandleFunc("/open", controller.HandleOpen)

	// Visualizer
	visualizer := NewVisualizer(debugger.tempDir.Dir)
	graph.ConnectTo(debugger.refReceiver.RefsOut, visualizer.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// Web UI
	webUpdater := NewWebUpdater(controller)
	graph.ConnectTo(visualizer.PathOut, webUpdater.PathIn)
	graph.ConnectTo(debugger.refReceiver.RefsOut, webUpdater.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// JSON combiner
	jsonCombiner := termites_web.NewJsonCombiner()
	graph.ConnectTo(visualizer.PathOut, jsonCombiner.JsonDataIn, termites.Via(VisualizerAdapter))
	graph.ConnectTo(jsonCombiner.JsonDataOut, connector.Hub.InFromApp)

	// Serve static files
	router.PathPrefix("/dbg-static/").Methods("GET").Handler(http.StripPrefix("/dbg-static/", http.FileServer(http.Dir(debugger.tempDir.Dir))))

	// Run termites_web server
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", debugger.httpPort), router); err != nil {
			panic(err)
		}
	}()
}

type Debugger struct {
	tempDir *termites.ManagedTempDirectory

	httpPort        int
	editor          CodeEditor
	refReceiver     *refReceiver
	messageReceiver *messageReceiver
}

type debuggerConfig struct {
	httpPort int
	editor   CodeEditor
}

// NewDebugger instantiates a non-connected Debugger, mainly available for advanced usage.
// Prefer to use the termites.GraphOption function WithDebugger to attach it directly to a graph if possible.
func NewDebugger(options ...DebuggerOption) *Debugger {
	config := &debuggerConfig{
		httpPort: 4242,
		editor:   nil,
	}

	for _, opt := range options {
		opt(config)
	}

	return &Debugger{
		tempDir: termites.NewManagedTempDirectory("debug-"),

		httpPort:        config.httpPort,
		editor:          config.editor,
		refReceiver:     newRefReceiver(),
		messageReceiver: newMsgReceiver(),
	}
}

func (d *Debugger) SetEventBus(b termites.EventBus) {
	b.Subscribe(termites.NodeRefUpdated, d.OnNodeRefUpdated)
	b.Subscribe(termites.NodeStopped, d.OnNodeStopped)
}

func (d *Debugger) OnNodeRefUpdated(e termites.Event) error {
	n, ok := e.Data.(termites.NodeUpdatedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	d.refReceiver.refChan <- n.Ref
	return nil
}

func (d *Debugger) OnNodeStopped(e termites.Event) error {
	n, ok := e.Data.(termites.NodeStoppedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	d.refReceiver.removeChan <- n.Id
	return nil
}
