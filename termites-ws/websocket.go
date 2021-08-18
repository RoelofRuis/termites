package cliserv

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO: fix appropriate message logging
type WebSocketClient struct {
	*Client

	conn *websocket.Conn

	readDeadline  time.Duration
	writeDeadline time.Duration
	pingInterval  time.Duration
}

type WebsocketConnector struct {
	registry ClientRegistry
}

func NewWebsocketConnector(registry ClientRegistry) *WebsocketConnector {
	return &WebsocketConnector{
		registry: registry,
	}
}

func (c *WebsocketConnector) ConnectWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	var id = ""
	keys, ok := r.URL.Query()["id"]
	if ok {
		id = keys[0]
	}

	client := c.registry.RegisterClient(id)

	wsclient := &WebSocketClient{
		Client:        client,
		conn:          conn,
		readDeadline:  60 * time.Second,
		writeDeadline: 10 * time.Second,
		pingInterval:  50 * time.Second,
	}

	go wsclient.writePump()
	go wsclient.readPump()
}

func (c *WebSocketClient) readPump() {
	defer func() {
		c.Unregister()
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	err := c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
	if err != nil {
		log.Printf("error setting read deadline: %v", err)
	}
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		c.Client.Send(message)
	}
}

func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(c.pingInterval)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Client.Received:
			err := c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
			if err != nil {
				log.Printf("error setting write deadline: %v", err)
			}
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("error writing message [%s]: %v", message, err)
			}

			if err := w.Close(); err != nil {
				log.Printf("error closing writer: %v", err)
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error writing ping message: %v", err)
				return
			}
		}
	}
}
