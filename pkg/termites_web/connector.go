package termites_web

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Connector struct {
	graph    *termites.Graph
	Hub      *Hub
	upgrader websocket.Upgrader

	clientIds map[string]bool
}

//go:embed connect.js
var embeddedJS embed.FS

func NewConnector(graph *termites.Graph, upgrader websocket.Upgrader) *Connector {
	return &Connector{
		graph:     graph,
		Hub:       newHub(),
		upgrader:  upgrader,
		clientIds: make(map[string]bool),
	}
}

func (c *Connector) Bind(router *mux.Router) {
	router.Path("/ws").Methods("GET").HandlerFunc(c.HandleWS)

	embeddedServer := http.FileServer(http.FS(embeddedJS))
	router.PathPrefix("/embedded/").Methods("GET").Handler(http.StripPrefix("/embedded/", embeddedServer))
}

func (c *Connector) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to upgrade connection: %s", err), http.StatusInternalServerError)
		return
	}

	var id = ""
	keys, ok := r.URL.Query()["id"]
	if ok {
		id = keys[0]
	}
	proposedId := id

	clientId := termites.RandomID()
	if c, has := c.clientIds[proposedId]; has && !c {
		clientId = proposedId
	}

	connectWebsocketIn(clientId, conn, c)
	connectWebSocketOut(clientId, conn, c)

	c.Hub.registerClient(clientId)
}
