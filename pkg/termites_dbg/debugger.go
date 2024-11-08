package termites_dbg

import (
	"encoding/json"
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

	// Messages
	if debugger.messageReceiver != nil {
		graph.Connect(debugger.messageReceiver.MessagesOut, connector.Hub.InFromApp, termites.Via(MessageSentAdapter))
	}

	// Logs
	if debugger.logReceiver != nil {
		graph.Connect(debugger.logReceiver.LogsOut, connector.Hub.InFromApp, termites.Via(LogsAdapter))
	}

	// Web UI
	router := mux.NewRouter()
	app, _ := App()
	router.Path("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(app))))
	router.Path("/ws").Methods("GET").Handler(connector)
	router.PathPrefix("/dbg-static/").Methods("GET").Handler(http.StripPrefix("/dbg-static/", http.FileServer(http.Dir(debugger.tempDir.Dir))))

	// State
	debuggerState, _ := json.Marshal(struct {
		GraphEnabled    bool `json:"graph_enabled"`
		MessagesEnabled bool `json:"messages_enabled"`
		LogsEnabled     bool `json:"logs_enabled"`
	}{
		GraphEnabled:    debugger.refReceiver != nil,
		MessagesEnabled: debugger.messageReceiver != nil,
		LogsEnabled:     debugger.logReceiver != nil,
	})
	initialState := make(map[string]json.RawMessage)
	initialState["debugger"] = debuggerState
	state := termites_web.NewStateWithInitial(initialState)

	graph.Connect(connector.Hub.ConnectionOut, state.ConnectionIn)
	graph.Connect(state.MessageOut, connector.Hub.InFromApp)

	// Visualizer
	if debugger.refReceiver != nil {
		visualizer := NewVisualizer(debugger.tempDir.Dir)
		graph.Connect(debugger.refReceiver.RefsOut, visualizer.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))
		graph.Connect(visualizer.PathOut, state.In, termites.Via(VisualizerAdapter))
	}

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
	refReceiver     *refReceiver
	messageReceiver *messageReceiver
	logReceiver     *logReceiver
}

type debuggerConfig struct {
	httpPort int
	upgrader websocket.Upgrader

	trackRefChanges bool
	trackMessages   bool
	trackLogs       bool
}

// NewDebugger instantiates a non-connected Debugger, mainly available for advanced usage.
// Prefer to use the termites.GraphOption function WithDebugger to attach it directly to a graph if possible.
func NewDebugger(options ...DebuggerOption) *Debugger {
	config := &debuggerConfig{
		httpPort: 4242,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		trackRefChanges: true,
		trackMessages:   true,
		trackLogs:       true,
	}

	for _, opt := range options {
		opt(config)
	}

	var r *refReceiver
	if config.trackRefChanges {
		r = newRefReceiver()
	}

	var m *messageReceiver
	if config.trackMessages {
		m = newMsgReceiver()
	}

	var l *logReceiver
	if config.trackLogs {
		l = newLogReceiver()
	}

	return &Debugger{
		tempDir: NewManagedTempDirectory("debug-"),

		httpPort:        config.httpPort,
		upgrader:        config.upgrader,
		refReceiver:     r,
		messageReceiver: m,
		logReceiver:     l,
	}
}

func (d *Debugger) SetEventBus(b termites.EventBus) {
	if d.refReceiver != nil {
		b.Subscribe(termites.NodeRefUpdated, d.refReceiver.onNodeRefUpdated)
		b.Subscribe(termites.NodeStopped, d.refReceiver.onNodeStopped)
	}

	if d.logReceiver != nil {
		b.Subscribe(termites.InfoLog, d.logReceiver.onLog)
		b.Subscribe(termites.ErrorLog, d.logReceiver.onLog)
		b.Subscribe(termites.PanicLog, d.logReceiver.onLog)
	}

	if d.messageReceiver != nil {
		b.Subscribe(termites.MessageSent, d.messageReceiver.onMessageSent)
	}
}
