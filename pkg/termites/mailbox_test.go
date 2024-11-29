package termites

import (
	"sync"
	"testing"
	"time"
)

func TestDefaultMailbox(t *testing.T) {
	n, mbox := configureNodeWithMailbox[int]()

	go func() {
		_ = mbox.deliver(Message{Data: 0})
	}()

	res, _ := n.ReceiveWithin(10 * time.Millisecond)

	if res != 0 {
		t.Errorf("received incorrect message")
	}
}

func TestMailboxWithTimeout(t *testing.T) {
	n, mbox := configureNodeWithMailbox[int](WithTimeout(10 * time.Millisecond))
	n.Delay = 100 * time.Millisecond

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = mbox.deliver(Message{Data: 0})
		err := mbox.deliver(Message{Data: 1})
		if err.Error() != "delivery timed out" {
			t.Errorf("expected message to time out")
		}
		wg.Done()
	}()

	res, _ := n.ReceiveWithin(1 * time.Second)
	if res != 0 {
		t.Errorf("received incorrect message")
	}

	wg.Wait()
}

func TestBufferedMailboxWithTimeout(t *testing.T) {
	n, mbox := configureNodeWithMailbox[int](WithBuffer(1), WithTimeout(10*time.Millisecond))
	n.Delay = 100 * time.Millisecond

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = mbox.deliver(Message{Data: 0})
		_ = mbox.deliver(Message{Data: 1})
		_ = mbox.deliver(Message{Data: 2})
		err := mbox.deliver(Message{Data: 3})
		if err.Error() != "delivery timed out" {
			t.Errorf("expected message to time out")
		}
		wg.Done()
	}()

	res, _ := n.ReceiveWithin(1 * time.Second)
	if res != 0 {
		t.Errorf("expected message %d, got %d", 0, res)
	}

	res, _ = n.ReceiveWithin(1 * time.Second)
	if res != 1 {
		t.Errorf("expected message %d, got %d", 1, res)
	}

	res, _ = n.ReceiveWithin(1 * time.Second)
	if res != 2 {
		t.Errorf("expected message %d, got %d", 2, res)
	}

	wg.Wait()
}

func TestDebouncedMailboxWithoutDelay(t *testing.T) {
	n, mbox := configureNodeWithMailbox[int](WithDebounce(100 * time.Millisecond))

	go func() {
		_ = mbox.deliver(Message{Data: 0})
		_ = mbox.deliver(Message{Data: 1})
	}()

	res, err := n.ReceiveWithin(1 * time.Second)
	if err != nil {
		t.Error(err)
	}

	if res != 1 {
		t.Errorf("received incorrect message")
	}
}

func TestDebouncedMailboxWithDelay(t *testing.T) {
	n, mbox := configureNodeWithMailbox[int](WithDebounce(100 * time.Millisecond))

	go func() {
		_ = mbox.deliver(Message{Data: 0})
		time.Sleep(200 * time.Millisecond)
		_ = mbox.deliver(Message{Data: 1})
	}()

	res, _ := n.ReceiveWithin(1 * time.Second)
	if res != 0 {
		t.Errorf("received incorrect message")
	}

	res, _ = n.ReceiveWithin(1 * time.Second)
	if res != 1 {
		t.Errorf("received incorrect message")
	}
}

func configureNodeWithMailbox[A any](opts ...MailboxOption) (*InspectableNode[A], *mailbox) {
	graph := NewGraph()
	n := NewInspectableNode[A]("test")
	mbox := newMailbox(n.In, opts...)

	conn := &Connection{
		id:      NewIdentifier("connection"),
		from:    nil,
		mailbox: mbox,
		adapter: nil,
	}

	mbox.to.connect(conn)
	graph.start(conn)

	return n, mbox
}
