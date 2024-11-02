package termites_web

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebSocketIn(t *testing.T) {
	graph := termites.NewGraph()
	connector := NewConnector(graph, upgrader)

	server := httptest.NewServer(connector)
	wsUrl := strings.ReplaceAll(server.URL, "http", "ws")
	wsConn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	ws := newWebsocketIn("test-in", wsConn)
	dataOut := termites.NewInspectableNode[ClientMessage]("ClientMessage")
	graph.Connect(ws.DataOut, dataOut.In)

	// Ping message
	if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
		t.Fatal(err)
	}
	tpe, data, err := wsConn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v %+v\n", tpe, data)
}
