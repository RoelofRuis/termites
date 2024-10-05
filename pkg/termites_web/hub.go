package termites_web

import (
	"time"

	"github.com/RoelofRuis/termites/pkg/termites"
)

type Hub struct {
	InFromWeb     *termites.InPort
	OutToApp      *termites.OutPort
	InFromApp     *termites.InPort
	OutToWeb      *termites.OutPort
	ConnectionOut *termites.OutPort
}

func newHub() *Hub {
	builder := termites.NewBuilder("Websocket Hub")

	h := &Hub{
		InFromWeb:     termites.NewInPort[ClientMessage](builder, "In From Web"),
		OutToApp:      termites.NewOutPort[ClientMessage](builder, "Out To App"),
		InFromApp:     termites.NewInPort[ClientMessage](builder, "In From App"),
		OutToWeb:      termites.NewOutPort[ClientMessage](builder, "Out To Web"),
		ConnectionOut: termites.NewOutPort[ClientConnection](builder, "Connection Info"),
	}

	builder.OnRun(h.Run)
	builder.OnShutdown(h.Shutdown)

	return h
}

func (h *Hub) registerClient(clientId string) {
	h.ConnectionOut.Send(ClientConnection{ConnType: ClientConnect, Id: clientId})
}

func (h *Hub) Run(c termites.NodeControl) error {
	var lastState []byte = nil

	for {
		select {
		case msg := <-h.InFromWeb.Receive():
			h.OutToApp.Send(msg.Data)

		case msg := <-h.InFromApp.Receive():
			clientMessage := msg.Data.(ClientMessage)

			var err error
			lastState, err = MakeUpdateMessage(clientMessage.Data)
			if err != nil {
				c.LogError("cannot send update message", err)
				continue
			}
			h.OutToWeb.Send(ClientMessage{ClientId: clientMessage.ClientId, Data: lastState})
		}
	}
}

func (h *Hub) Shutdown(_ termites.TeardownControl) error {
	msg, err := MakeCloseMessage()
	if err != nil {
		return err
	}
	h.OutToWeb.Send(ClientMessage{Data: msg})
	time.Sleep(50 * time.Millisecond) // Give time to send close message
	return nil
}
