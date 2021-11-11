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

const WebURL = "localhost:8000"

func main() {
	graph := termites.NewGraph(
		termites.PrintLogsToConsole(),
		termites.PrintMessagesToConsole(),
		termites_dbg.WithDebugger(),
	)

	connector := termites_web.NewConnector(graph)
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handleIndex)
	connector.Bind(router)

	generator := examples.NewGenerator(100 * time.Millisecond)
	graph.ConnectTo(generator.BytesOut, connector.Hub.InFromApp)

	go func() {
		err := http.ListenAndServe(WebURL, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	_ = termites_web.RunBrowser(WebURL)

	graph.Wait()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./examples/websocket/index.html")
}
