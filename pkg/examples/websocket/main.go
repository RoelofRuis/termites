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
		termites.PrintMessagesToConsole(),
		termites_dbg.WithDebugger(),
	)

	// Create a router and bind the appropriate handlers
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handleIndex)

	// Create a new web connector
	connector := termites_web.NewConnector(graph)
	connector.Bind(router)

	// For demo purposes we create a generator.
	// This is where the custom application logic would go.
	generator := examples.NewGenerator(100 * time.Millisecond)
	graph.ConnectTo(generator.ClientOut, connector.Hub.InFromApp)

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
