package termites_web

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TestConnector_ConnectWebsocket(t *testing.T) {
	graph := termites.NewGraph()

	connector := NewConnector(graph, upgrader)
	clientIn := termites.NewInspectableNode[ClientConnection]("ClientConnection")
	messageIn := termites.NewInspectableNode[ClientMessage]("ClientMessage")

	graph.ConnectTo(connector.Hub.ConnectionOut, clientIn.In)
	graph.ConnectTo(connector.Hub.OutToApp, messageIn.In)

	server := httptest.NewServer(connector)
	wsUrl := strings.ReplaceAll(server.URL, "http", "ws")

	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteMessage(websocket.BinaryMessage, []byte("hello world")); err != nil {
		t.Fatal(err)
	}

	_, err = clientIn.ReceiveWithin(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := messageIn.ReceiveWithin(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if string(msg.Data) != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", string(msg.Data))
	}
}
