package termites_ws2

import (
	"github.com/RoelofRuis/termites/termites"
	"time"
)

type Hub struct {
	InFromWeb     *termites.InPort
	OutToApp      *termites.OutPort
	InFromApp     *termites.InPort
	OutToWeb      *termites.OutPort
	ConnectionOut *termites.OutPort
}

func NewHub() *Hub {
	builder := termites.NewBuilder("Websocket Hub")

	h := &Hub{
		InFromWeb:     builder.InPort("In From Web", ClientMessage{}),
		OutToApp:      builder.OutPort("Out To App", ClientMessage{}),
		InFromApp:     builder.InPort("In From App", []byte{}),
		OutToWeb:      builder.OutPort("Out To Web", []byte{}),
		ConnectionOut: builder.OutPort("Client status out", ClientConnection{}),
	}

	builder.OnRun(h.Run)
	builder.OnShutdown(h.Shutdown)

	return h
}

func (h *Hub) registerClient(client Client) {
	// register and send on ConnectionOut
}

func (h *Hub) Run(c termites.NodeControl) error {
	var lastState []byte = nil

	for {
		select {
		case msg := <-h.InFromWeb.Receive():
			h.OutToApp.Send(msg.Data)

		case msg := <-h.InFromApp.Receive():
			var err error
			lastState, err = MakeUpdateMessage(msg.Data.([]byte))
			if err != nil {
				c.LogError("cannot send update message", err)
				continue
			}
			h.OutToWeb.Send(lastState)
		}
	}
}

func (h *Hub) Shutdown(_ termites.TeardownControl) error {
	msg, err := MakeCloseMessage()
	if err != nil {
		return err
	}
	h.OutToWeb.Send(msg)
	time.Sleep(50 * time.Millisecond) // Give time to send close message
	return nil
}
