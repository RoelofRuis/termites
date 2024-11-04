package termites_web

import (
	"fmt"
	"net/http"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/websocket"
)

type Connector struct {
	graph    *termites.Graph
	Hub      *Hub
	upgrader websocket.Upgrader

	clientIds map[string]bool
}

func NewConnector(graph *termites.Graph, upgrader websocket.Upgrader) *Connector {
	return &Connector{
		graph:     graph,
		Hub:       newHub(),
		upgrader:  upgrader,
		clientIds: make(map[string]bool),
	}
}

func (c *Connector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	wsIn := newWebsocketIn(clientId, conn)
	wsIn.graphConnection = c.graph.Connect(wsIn.DataOut, c.Hub.InFromWeb)

	wsOut := newWebsocketOut(clientId, conn)
	wsOut.graphConnection = c.graph.Connect(c.Hub.OutToWeb, wsOut.DataIn)

	c.Hub.registerClient(clientId)
}
