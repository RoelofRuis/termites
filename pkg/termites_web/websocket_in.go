package termites_web

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/gorilla/websocket"
	"time"
)

type webSocketIn struct {
	DataOut *termites.OutPort

	id              string
	conn            *websocket.Conn
	readDeadline    time.Duration
	readLimit       int64
	graphConnection *termites.Connection
}

func newWebsocketIn(id string, conn *websocket.Conn, readLimit int64) *webSocketIn {
	builder := termites.NewBuilder("websocket IN")

	ws := &webSocketIn{
		DataOut: termites.NewOutPort[ClientMessage](builder),

		id:           id,
		conn:         conn,
		readDeadline: 60 * time.Second,
		readLimit:    readLimit,
	}

	builder.OnRun(ws.Run)

	return ws
}

func (w *webSocketIn) Run(c termites.NodeControl) error {
	defer func() {
		_ = w.conn.Close()
		w.graphConnection.Disconnect()
	}()

	w.conn.SetReadLimit(w.readLimit)
	if err := w.conn.SetReadDeadline(time.Now().Add(w.readDeadline)); err != nil {
		return err
	}

	w.conn.SetPongHandler(func(string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(w.readDeadline))
		return nil
	})

	for {
		if err := w.conn.SetReadDeadline(time.Now().Add(w.readDeadline)); err != nil {
			return err
		}

		_, message, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				c.LogError("websocket unexpected close error: %s", err)
			}
			break
		}
		w.DataOut.Send(ClientMessage{ClientId: w.id, Data: message})
	}

	return nil
}
