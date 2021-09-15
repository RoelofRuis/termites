package termites_dbg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RoelofRuis/termites/termites_web"
	"github.com/gorilla/mux"

	"github.com/RoelofRuis/termites/termites"
)

func WithDebugger(httpPort int) termites.GraphOptions {
	graph := termites.NewGraph(
		termites.Named("Termites Debugger"),
		termites.WithoutSigtermHandler(),
	)

	dbg := NewDebugger(httpPort)

	Init(graph, dbg)

	return termites.WithEventSubscriber(dbg)
}

func Init(graph termites.Graph, debugger *debugger) {
	connector := termites_web.NewConnector(graph)

	router := mux.NewRouter()
	connector.Bind(router)

	// Visualizer
	visualizer := NewVisualizer()
	graph.ConnectTo(debugger.refReceiver.RefsOut, visualizer.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// Web UI
	webUI := NewWebController(NewWebUI(router), debugger.staticDir)
	graph.ConnectTo(visualizer.PathOut, webUI.PathIn)
	graph.ConnectTo(debugger.refReceiver.RefsOut, webUI.RefsIn, termites.WithMailbox(&termites.DebouncedMailbox{Delay: 100 * time.Millisecond}))

	// JSON combiner
	jsonCombiner := termites_web.NewJsonCombiner()
	graph.ConnectTo(visualizer.PathOut, jsonCombiner.JsonDataIn, termites.Via(VisualizerAdapter))
	graph.ConnectTo(jsonCombiner.JsonDataOut, connector.Hub.InFromApp)

	// Serve static files
	router.PathPrefix("/dbg-static/").Methods("GET").Handler(http.StripPrefix("/dbg-static/", http.FileServer(http.Dir(debugger.staticDir))))

	// Run web server
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", debugger.httpPort), router); err != nil {
			panic(err)
		}
	}()
}

type debugger struct {
	httpPort        int
	staticDir       string
	refReceiver     *refReceiver
	messageReceiver *messageReceiver
}

// NewDebugger instantiates a non-connected debugger, mainly available for advanced usage.
// Prefer to use the termites.GraphOptions function WithDebugger to attach it directly to a graph if possible.
func NewDebugger(httpPort int) *debugger {
	staticDir, err := ioutil.TempDir("", "debug-")
	if err != nil {
		panic(err)
	}

	return &debugger{
		httpPort:        httpPort,
		staticDir:       staticDir,
		refReceiver:     newRefReceiver(),
		messageReceiver: newMsgReceiver(),
	}
}

func (d *debugger) SetEventBus(b termites.EventBus) {
	b.Subscribe(termites.NodeRefUpdated, d.OnNodeRefUpdated)
	b.Subscribe(termites.NodeStopped, d.OnNodeStopped)
}

func (d *debugger) OnNodeRefUpdated(e termites.Event) error {
	n, ok := e.Data.(termites.NodeUpdatedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	d.refReceiver.refChan <- n.Ref
	return nil
}

func (d *debugger) OnNodeStopped(e termites.Event) error {
	n, ok := e.Data.(termites.NodeStoppedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	d.refReceiver.removeChan <- n.Id
	return nil
}
