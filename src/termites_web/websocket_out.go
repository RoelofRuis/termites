package termites_web

import (
	"fmt"
	"github.com/RoelofRuis/termites/termites"
	"github.com/gorilla/websocket"
	"time"
)

type webSocketOut struct {
	DataIn *termites.InPort

	id              string
	conn            *websocket.Conn
	pingInterval    time.Duration
	writeDeadline   time.Duration
	graphConnection *termites.Connection
}

func connectWebSocketOut(id string, conn *websocket.Conn, connector *connector) {
	builder := termites.NewBuilder("websocket OUT")

	ws := &webSocketOut{
		DataIn: builder.InPort("Data In", []byte{}),

		id:            id,
		conn:          conn,
		pingInterval:  50 * time.Second,
		writeDeadline: 10 * time.Second,
	}

	builder.OnRun(ws.Run)

	outConn := connector.graph.ConnectTo(connector.Hub.OutToWeb, ws.DataIn)
	ws.graphConnection = outConn
}

func (w *webSocketOut) Run(c termites.NodeControl) error {
	ticker := time.NewTicker(w.pingInterval)
	defer func() {
		ticker.Stop()
		_ = w.conn.Close()
		w.graphConnection.Disconnect()
	}()

	for {
		select {
		case message, ok := <-w.DataIn.Receive():
			bytes := message.Data.([]byte)
			if err := w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline)); err != nil {
				c.LogError("error setting write deadline", err)
			}
			if !ok {
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}

			w, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return nil
			}
			if _, err = w.Write(bytes); err != nil {
				return fmt.Errorf("write failed: %w", err)
			}
			if err = w.Close(); err != nil {
				return fmt.Errorf("unable to close writer: %w", err)
			}

		case <-ticker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("unable to send ping message: %w", err)
			}
		}
	}
}
