package main

import (
	"github.com/RoelofRuis/termites/examples"
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
	"github.com/RoelofRuis/termites/termites_ws2"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	graph := termites.NewGraph(termites.WithConsoleLogger(), termites_dbg.WithDebugger(4242))

	router := mux.NewRouter()
	connector := termites_ws2.NewConnector(graph)
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

	graph.Wait()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./examples/websocket/static/index.html")
}
