package termites_web

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
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

func newWebsocketOut(id string, conn *websocket.Conn) *webSocketOut {
	builder := termites.NewBuilder("websocket OUT")

	ws := &webSocketOut{
		DataIn: termites.NewInPort[ClientMessage](builder),

		id:            id,
		conn:          conn,
		pingInterval:  50 * time.Second,
		writeDeadline: 10 * time.Second,
	}

	builder.OnRun(ws.Run)

	return ws
}

func (w *webSocketOut) Run(c termites.NodeControl) error {
	ticker := time.NewTicker(w.pingInterval)
	defer func() {
		ticker.Stop()
		_ = w.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
		_ = w.conn.Close()
		w.graphConnection.Disconnect()
	}()

	for {
		select {
		case msg, ok := <-w.DataIn.Receive():
			if !ok {
				return nil
			}

			clientMessage := msg.Data.(ClientMessage)
			if clientMessage.ClientId != "" && clientMessage.ClientId != w.id {
				continue // drop message not intended for this client
			}

			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return nil
			}

			if err := w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline)); err != nil {
				c.LogError("error setting write deadline", err)
			}

			if _, err = writer.Write(clientMessage.Data); err != nil {
				return fmt.Errorf("write failed: %w", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("unable to close writer: %w", err)
			}

		case <-ticker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline))
			if err := w.conn.WriteControl(websocket.PingMessage, []byte{}, time.Time{}); err != nil {
				return fmt.Errorf("unable to send ping message: %w", err)
			}
		}
	}
}
