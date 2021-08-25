package web

import (
	"github.com/RoelofRuis/termites/termites-core"
	cliserv "github.com/RoelofRuis/termites/termites-ws"
	"log"
)

type In struct {
	ConnectionOut *termites.OutPort

	hub cliserv.Hub
}

func NewIn(hub cliserv.Hub) *In {
	builder := termites.NewBuilder("Web In")

	in := &In{
		ConnectionOut: builder.OutPort("Connection", cliserv.ClientConnection{}),
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
