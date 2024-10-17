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

type ClientConnection struct {
	ConnType ConnectionType
	Id       string
}

type ConnectionType uint8

const (
	ClientConnect ConnectionType = 0
)

func newHub() *Hub {
	builder := termites.NewBuilder("Websocket Hub")

	h := &Hub{
		InFromWeb:     termites.NewInPortNamed[ClientMessage](builder, "In From Web"),
		OutToApp:      termites.NewOutPortNamed[ClientMessage](builder, "Out To App"),
		InFromApp:     termites.NewInPortNamed[ClientMessage](builder, "In From App"),
		OutToWeb:      termites.NewOutPortNamed[ClientMessage](builder, "Out To Web"),
		ConnectionOut: termites.NewOutPort[ClientConnection](builder),
	}

	builder.OnRun(h.Run)
	builder.OnShutdown(h.Shutdown)

	return h
}

func (h *Hub) registerClient(clientId string) {
	h.ConnectionOut.Send(ClientConnection{ConnType: ClientConnect, Id: clientId})
}

func (h *Hub) Run(_ termites.NodeControl) error {
	for {
		select {
		case msg := <-h.InFromWeb.Receive():
			h.OutToApp.Send(msg.Data)

		case msg := <-h.InFromApp.Receive():
			clientMessage := msg.Data.(ClientMessage)
			h.OutToWeb.Send(clientMessage)
		}
	}
}

func (h *Hub) Shutdown(_ termites.TeardownControl) error {
	msg, err := NewClientMessage(SystemCloseTopic, nil)
	if err != nil {
		return err
	}
	h.OutToWeb.Send(msg)
	time.Sleep(50 * time.Millisecond) // Give time to send close message
	return nil
}
