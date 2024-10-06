package termites_dbg

import (
	"fmt"
	"github.com/gorilla/websocket"
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
	connector := termites_web.NewConnector(graph, debugger.upgrader)

	router := mux.NewRouter()
	connector.Bind(router)

	controller := NewWebController()
	controller.editor = debugger.editor

	router.HandleFunc("/", controller.HandleIndex)
	router.HandleFunc("/open", controller.HandleOpen)

	// Visualizer
	visualizer := NewVisualizer(debugger.tempDir.Dir)
	graph.ConnectTo(debugger.refReceiver.RefsOut, visualizer.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// Open References
	webUpdater := NewWebUpdater(controller)
	graph.ConnectTo(debugger.refReceiver.RefsOut, webUpdater.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// State
	stateTracker := termites_web.NewStateTracker()
	graph.ConnectTo(connector.Hub.ConnectionOut, stateTracker.ConnectionIn)
	graph.ConnectTo(stateTracker.MessageOut, connector.Hub.InFromApp)

	graph.ConnectTo(visualizer.PathOut, stateTracker.StateIn, termites.Via("Visualizer adapter", VisualizerAdapter))

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
	tempDir *ManagedTempDirectory

	httpPort        int
	upgrader        websocket.Upgrader
	editor          CodeEditor
	refReceiver     *refReceiver
	messageReceiver *messageReceiver
}

type debuggerConfig struct {
	httpPort int
	editor   CodeEditor
	upgrader websocket.Upgrader
}

// NewDebugger instantiates a non-connected Debugger, mainly available for advanced usage.
// Prefer to use the termites.GraphOption function WithDebugger to attach it directly to a graph if possible.
func NewDebugger(options ...DebuggerOption) *Debugger {
	config := &debuggerConfig{
		httpPort: 4242,
		editor:   nil,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	for _, opt := range options {
		opt(config)
	}

	return &Debugger{
		tempDir: NewManagedTempDirectory("debug-"),

		httpPort:        config.httpPort,
		upgrader:        config.upgrader,
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
