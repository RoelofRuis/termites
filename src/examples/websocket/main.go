package main

import (
	"github.com/RoelofRuis/termites/examples"
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
	"github.com/RoelofRuis/termites/termites_web"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	graph := termites.NewGraph(
		termites.PrintLogsToConsole(),
		termites.PrintMessagesToConsole(),
		termites_dbg.WithDebugger(4242),
	)

	connector := termites_web.NewConnector(graph)

	router := mux.NewRouter()
	connector.Bind(router)

	generator := examples.NewGenerator(100 * time.Millisecond)
	graph.ConnectTo(generator.BytesOut, connector.Hub.InFromApp)

	router.Path("/").Methods("GET").HandlerFunc(handleIndex)

	go func() {
		err := http.ListenAndServe(":8000", router)
		if err != nil {
			log.Println(err)
		}
	}()

	_ = termites_web.RunBrowser("localhost:8000")

	graph.Wait()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./examples/websocket/static/index.html")
}
