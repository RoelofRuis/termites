package main

import (
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

// Setting up a full duplex websocket connection
func main() {
	// Create a new graph
	graph := termites.NewGraph(
		termites.PrintLogsToConsole(),
		//termites.PrintMessagesToConsole(),
		termites_dbg.WithDebugger(), // Visit http://localhost:4242 to view the debugger
	)

	// Create a new web connector
	connector := termites_web.NewConnector(graph)

	// Create a router and bind the appropriate handlers
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handleIndex)
	router.Path("/ws").Methods("GET").Handler(connector)
	router.PathPrefix("/embedded/").Methods("GET").Handler(http.StripPrefix("/embedded/", termites_web.EmbeddedJS()))

	// For demo purposes we create a generator.
	// This is where the custom application logic would go.
	generator := examples.NewGenerator(1000 * time.Millisecond)

	// We collect the web-sharable state in a state instance
	state := termites_web.NewState()
	graph.Connect(state.MessageOut, connector.Hub.InFromApp)
	graph.Connect(connector.Hub.ConnectionOut, state.ConnectionIn)
	graph.Connect(generator.CountOut, state.In, termites.Via(termites_web.MarshalState("generator")))

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
