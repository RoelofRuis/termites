package termites_dbg

import (
	"github.com/RoelofRuis/termites/pkg/termites"
)

// messageReceiver handles incoming messages from the target graphs event bus.
type messageReceiver struct {
	MessagesOut *termites.OutPort
}

func newMsgReceiver() *messageReceiver {
	builder := termites.NewBuilder("Msg Receiver")

	n := &messageReceiver{
		MessagesOut: termites.NewOutPortNamed[termites.MessageSentEvent](builder, "Events"),
	}

	return n
}

func (r *messageReceiver) onMessageSent(e termites.Event) error {
	msg, ok := e.Data.(termites.MessageSentEvent)
	if !ok {
		return termites.InvalidEventError
	}
	r.MessagesOut.Send(msg)
	return nil
}
