package termites

import (
	"reflect"
	"testing"
	"time"
)

func TestDebouncingMailboxWithoutDelay(t *testing.T) {
	port := newInPort("test", reflect.TypeOf(""), nil)
	mailbox := mailboxFromConfig(port, &DebouncedMailbox{Delay: 100 * time.Millisecond})

	go func() {
		_ = mailbox.deliver(Message{Data: 0})
		_ = mailbox.deliver(Message{Data: 1})
	}()

	res := <-port.Receive()

	if res.Data != 1 {
		t.Errorf("received incorrect message")
	}
}

func TestDebouncingMailboxWithDelay(t *testing.T) {
	port := newInPort("test", reflect.TypeOf(""), nil)
	mailbox := mailboxFromConfig(port, &DebouncedMailbox{Delay: 100 * time.Millisecond})

	go func() {
		_ = mailbox.deliver(Message{Data: 0})
		time.Sleep(200 * time.Millisecond)
		_ = mailbox.deliver(Message{Data: 1})
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
