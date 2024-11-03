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

func TestWebSocketOut(t *testing.T) {
	graph := termites.NewGraph()
	dataIn := termites.NewInspectableNode[ClientMessage]("ClientMessage")

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		ws := newWebsocketOut("test-out", conn)
		graph.Connect(dataIn.Out, ws.DataIn)
	})

	server := httptest.NewServer(router)
	wsUrl := strings.ReplaceAll(server.URL, "http", "ws")

	wsConn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	dataIn.Send <- ClientMessage{Data: []byte("Test 123")}

	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	if string(msg) != "Test 123" {
		t.Errorf("Message was incorrect, got: %s, want: %s.", msg, "Test 123")
	}

	if err := wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{}); err != nil {
		t.Fatal(err)
	}

	if err := wsConn.Close(); err != nil {
		t.Fatal(err)
	}

	dataIn.Send <- ClientMessage{Data: []byte("Test 123")}

	_, msg, err = wsConn.ReadMessage()
	if err == nil {
		t.Fatal("expected tcp connection to be closed")
	}
}
