package web

import (
	"log"

	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_ws"
)

type In struct {
	ConnectionOut *termites.OutPort

	hub termites_ws.Hub
}

func NewIn(hub termites_ws.Hub) *In {
	builder := termites.NewBuilder("Web In")

	in := &In{
		ConnectionOut: builder.OutPort("Connection", termites_ws.ClientConnection{}),
		hub:           hub,
	}

	builder.OnRun(in.Run)

	return in
}

func (w *In) Run(_ termites.NodeControl) error {
	for {
		select {
		case msg := <-w.hub.ReadConnect():
			log.Printf("WEB CONNECT: %+v\n", msg)
			w.ConnectionOut.Send(msg)

		case msg := <-w.hub.ReadReceive():
			log.Printf("WEB RECEIVE: %+v\n", msg)
		}
	}
}
