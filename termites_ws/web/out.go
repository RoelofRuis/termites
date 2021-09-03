package web

import (
	"log"
	"time"

	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_ws"
)

type Out struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort

	hub termites_ws.Hub
}

func NewOut(hub termites_ws.Hub) *Out {
	builder := termites.NewBuilder("Web Out")

	out := &Out{
		ConnectionIn: builder.InPort("Connection", termites_ws.ClientConnection{}),
		StateIn:      builder.InPort("State", []byte{}),
		hub:          hub,
	}

	builder.OnRun(out.Run)
	builder.OnShutdown(out.Shutdown)

	return out
}

func (w *Out) Run(_ termites.NodeControl) error {
	var lastState []byte = nil

	for {
		select {
		case msg := <-w.StateIn.Receive():
			var err error
			lastState, err = termites_ws.MakeUpdateMessage(msg.Data.([]byte))
			if err != nil {
				log.Printf("Web Out: cannot send update message: %s", err.Error())
				continue
			}
			w.hub.Broadcast(lastState)

		case <-w.ConnectionIn.Receive():
			w.hub.Broadcast(lastState)
		}
	}
}

func (w *Out) Shutdown(_ time.Duration) error {
	msg, err := termites_ws.MakeCloseMessage()
	if err != nil {
		return err
	}
	w.hub.Broadcast(msg)
	time.Sleep(50 * time.Millisecond) // Give time to send close message
	return nil
}
