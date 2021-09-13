package termites_dbg

import "github.com/RoelofRuis/termites/termites"

type messageReceiver struct {
	MessagesOut *termites.OutPort

	messageChan chan termites.MessageSentEvent
}

func newMsgReceiver() *messageReceiver {
	builder := termites.NewBuilder("Msg Receiver")

	n := &messageReceiver{
		MessagesOut: builder.OutPort("Messages", termites.MessageSentEvent{}),
	}

	builder.OnRun(n.run)

	return n
}

func (r *messageReceiver) run(_ termites.NodeControl) error {
	for msg := range r.messageChan {
		r.MessagesOut.Send(msg)
	}
	return nil
}
