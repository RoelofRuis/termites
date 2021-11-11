package termites_web

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connector struct {
	graph termites.Graph
	Hub   *Hub

	clientIds map[string]bool
}

//go:embed connect.js
var embeddedJS embed.FS

func NewConnector(graph termites.Graph) *connector {
	return &connector{
		graph: graph,
		Hub:   newHub(),

		clientIds: make(map[string]bool),
	}
}

func (c *connector) Bind(router *mux.Router) {
	router.Path("/ws").Methods("GET").HandlerFunc(c.ConnectWebsocket)

	embeddedServer := http.FileServer(http.FS(embeddedJS))
	router.PathPrefix("/embedded/").Methods("GET").Handler(http.StripPrefix("/embedded/", embeddedServer))
}

func (c *connector) ConnectWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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
