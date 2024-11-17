package termites_web

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TestWebSocketIn(t *testing.T) {
	graph := termites.NewGraph()
	dataOut := termites.NewInspectableNode[ClientMessage]("ClientMessage")

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		ws := newWebsocketIn("test-in", conn, 512)
		graph.Connect(ws.DataOut, dataOut.In)
	})

	server := httptest.NewServer(router)
	wsUrl := strings.ReplaceAll(server.URL, "http", "ws")

	wsConn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := wsConn.WriteMessage(websocket.TextMessage, []byte("Test 123")); err != nil {
		t.Fatal(err)
	}

	message, err := dataOut.ReceiveWithin(1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if string(message.Data) != "Test 123" {
		t.Errorf("received incorrect message")
	}

	if err := wsConn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
}
