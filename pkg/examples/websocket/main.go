package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"

	"github.com/RoelofRuis/termites/pkg/examples"
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_dbg"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"github.com/gorilla/mux"
)

const WebURL = "localhost:8008"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Setting up a full duplex websocket connection
func main() {
	// Create a new graph
	graph := termites.NewGraph(
		termites.PrintLogsToConsole(),
		termites.PrintMessagesToConsole(),
		termites_dbg.WithDebugger(), // Visit http://localhost:4242 to view the debugger
	)

	// Create a router and bind the appropriate handlers
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handleIndex)

	// Create a new web connector
	connector := termites_web.NewConnector(graph, upgrader)
	connector.Bind(router)

	// For demo purposes we create a generator.
	// This is where the custom application logic would go.
	generator := examples.NewGenerator(100 * time.Millisecond)

	// We create a state tracker that collects the web-sharable state
	stateTracker := termites_web.NewStateTracker()
	graph.ConnectTo(stateTracker.MessageOut, connector.Hub.InFromApp)
	graph.ConnectTo(connector.Hub.ConnectionOut, stateTracker.ConnectionIn)
	graph.ConnectTo(generator.IntOut, stateTracker.StateIn, termites.Via("generator adapter", generatorAdapter))

	go func() {
		// Run the webserver
		err := http.ListenAndServe(WebURL, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Open a browser
	_ = termites_web.RunBrowser(WebURL)

	// Wait for graph termination
	graph.Wait()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./examples/websocket/index.html")
}

func generatorAdapter(i int) (termites_web.StateMessage, error) {
	msg, err := json.Marshal(struct {
		Count int `json:"count"`
	}{
		Count: i,
	})
	if err != nil {
		return termites_web.StateMessage{}, err
	}

	return termites_web.StateMessage{Key: "generator", Data: msg}, nil
}
