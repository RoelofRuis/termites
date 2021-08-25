package web

import (
	"github.com/RoelofRuis/termites/termites-core"
	cliserv "github.com/RoelofRuis/termites/termites-ws"
	"log"
	"time"
)

type Out struct {
	ConnectionIn *termites.InPort
	StateIn      *termites.InPort

	hub cliserv.Hub
}

func NewOut(hub cliserv.Hub) *Out {
	builder := termites.NewBuilder("Web Out")

	out := &Out{
		ConnectionIn: builder.InPort("Connection", cliserv.ClientConnection{}),
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
			lastState, err = cliserv.MakeUpdateMessage(msg.Data.([]byte))
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
	msg, err := cliserv.MakeCloseMessage()
	if err != nil {
		return err
	}
	w.hub.Broadcast(msg)
	time.Sleep(50 * time.Millisecond) // Give time to send close message
	return nil
}
