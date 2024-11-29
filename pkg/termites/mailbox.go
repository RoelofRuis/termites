package termites

import (
	"time"
)

type mailbox struct {
	to      *InPort
	deliver func(msg Message) error
}

func newMailbox(to *InPort, opts ...MailboxOption) *mailbox {
	config := &mailboxConfig{
		capacity:         0,
		receiveTimeout:   0,
		debounceDelay:    0,
		errorWhenDropped: nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	receiver := to.receive

	if config.capacity > 0 {
		buffer := make(chan Message, config.capacity)

		go send(buffer, receiver)

		receiver = buffer
	}

	if config.debounceDelay > 0 {
		debounceChan := make(chan Message)

		go sendDebounced(debounceChan, receiver, config.debounceDelay)

		return &mailbox{to: to, deliver: deliverBlocking(debounceChan)}
	} else {
		return &mailbox{to: to, deliver: config.buildDeliverFunc(receiver)}
	}
}

func (c *mailboxConfig) buildDeliverFunc(to chan Message) func(msg Message) error {
	if c.receiveTimeout > 0 {
		return deliverWithTimeout(to, c.receiveTimeout, c.errorWhenDropped)
	}
	if c.receiveTimeout == 0 {
		return deliverBlocking(to)
	}
	return deliverNonBlocking(to, c.errorWhenDropped)
}

func deliverNonBlocking(receive chan Message, errorWhenDropped error) func(msg Message) error {
	return func(msg Message) error {
		select {
		case receive <- msg:
			return nil
		default:
			return errorWhenDropped
		}
	}
}

func deliverBlocking(receive chan Message) func(msg Message) error {
	return func(msg Message) error {
		receive <- msg
		return nil
	}
}

func deliverWithTimeout(receive chan Message, timeout time.Duration, errorWhenDropped error) func(msg Message) error {
	return func(msg Message) error {
		select {
		case receive <- msg:
			return nil
		case <-time.After(timeout):
			return errorWhenDropped
		}
	}
}

func send(from chan Message, to chan Message) {
	for msg := range from {
		to <- msg
	}
}

func sendDebounced(from chan Message, to chan Message, delay time.Duration) {
	var lastMessage Message
	for msg := range from {
		lastMessage = msg
	nextMessage:
		select {
		case <-time.After(delay):
			// TODO: this can still block, not sure whether to do something about it...
			to <- lastMessage
		case nextMsg := <-from:
			lastMessage = nextMsg
			goto nextMessage
		}
	}
}
