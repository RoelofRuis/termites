package termites_ws2

import (
	"github.com/RoelofRuis/termites/termites"
	"github.com/gorilla/websocket"
	"time"
)

type webSocketOut struct {
	DataIn *termites.InPort

	id            string
	conn          *websocket.Conn
	pingInterval  time.Duration
	writeDeadline time.Duration
}

func newWebSocketOut(id string, conn *websocket.Conn) *webSocketOut {
	builder := termites.NewBuilder("websocket OUT")

	ws := &webSocketOut{
		DataIn: builder.InPort("Data In", []byte{}),

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
		_ = w.conn.Close()
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
				return err
			}
			if _, err = w.Write(bytes); err != nil {
				return err
			}
			if err := w.Close(); err != nil {
				return err
			}

		case <-ticker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}
