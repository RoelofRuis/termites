package termites_web

import (
	"fmt"
	"net/http"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/websocket"
)

type Connector struct {
	graph     *termites.Graph
	Hub       *Hub
	upgrader  websocket.Upgrader
	readLimit int64

	clientIds map[string]bool
}

type connectorConfig struct {
	upgrader websocket.Upgrader
	// Websocket read limit in bytes
	readLimit int64
}

func NewConnector(graph *termites.Graph, options ...ConnectorOption) *Connector {
	config := &connectorConfig{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		readLimit: 1024,
	}

	for _, opt := range options {
		opt(config)
	}

	return &Connector{
		graph:     graph,
		Hub:       newHub(),
		upgrader:  config.upgrader,
		readLimit: config.readLimit,
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

	wsIn := newWebsocketIn(clientId, conn, 10000)
	wsIn.graphConnection = c.graph.Connect(wsIn.DataOut, c.Hub.InFromWeb)

	wsOut := newWebsocketOut(clientId, conn)
	wsOut.graphConnection = c.graph.Connect(c.Hub.OutToWeb, wsOut.DataIn)

	c.Hub.registerClient(clientId)
}
