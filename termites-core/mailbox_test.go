package termites

import (
	"testing"
	"time"
)

var factory = MailboxFactory{}

func TestDebouncingMailboxWithoutDelay(t *testing.T) {
	port := newInPort("test", "", nil)
	mailbox := factory.FromConfig(port, &DebouncedMailbox{Delay: 100 * time.Millisecond})

	go func() {
		mailbox.deliver(Message{Data: 0})
		mailbox.deliver(Message{Data: 1})
	}()

	res := <-port.Receive()

	if res.Data != 1 {
		t.Errorf("received incorrect message")
	}
}

func TestDebouncingMailboxWithDelay(t *testing.T) {
	port := newInPort("test", "", nil)
	mailbox := factory.FromConfig(port, &DebouncedMailbox{Delay: 100 * time.Millisecond})

	go func() {
		mailbox.deliver(Message{Data: 0})
		time.Sleep(200 * time.Millisecond)
		mailbox.deliver(Message{Data: 1})
	}()

	res := <-port.Receive()
	if res.Data != 0 {
		t.Errorf("received incorrect message")
	}

	res = <-port.Receive()
	if res.Data != 1 {
		t.Errorf("received incorrect message")
	}
}
