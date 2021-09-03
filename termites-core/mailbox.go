package termites

import (
	"time"
)

type Message struct {
	Data interface{}
}

type MailboxConfig interface {
	IsMailboxConfig()
}

type NormalMailbox struct {
	ReceiveTimeout time.Duration
}

func (m *NormalMailbox) IsMailboxConfig() {}

type CapacityMailbox struct {
	Capacity       int
	ReceiveTimeout time.Duration
}

func (m *CapacityMailbox) IsMailboxConfig() {}

type DebouncedMailbox struct {
	Delay time.Duration
}

func (m *DebouncedMailbox) IsMailboxConfig() {}

func mailboxFromConfig(to *InPort, c MailboxConfig) *mailbox {
	var deliverFunc func(msg Message) bool
	switch conf := c.(type) {
	case *NormalMailbox:
		deliverFunc = func(msg Message) bool {
			ticker := time.NewTimer(conf.ReceiveTimeout)
			select {
			case <-ticker.C:
				return false

			case to.receive <- msg:
				return true
			}
		}

	case *CapacityMailbox:
		messages := make(chan Message, conf.Capacity)

		go func() {
			for msg := range messages {
				to.receive <- msg
			}
		}()

		deliverFunc = func(msg Message) bool {
			ticker := time.NewTimer(conf.ReceiveTimeout)
			select {
			case <-ticker.C:
				return false

			case messages <- msg:
				return true
			}
		}

	case *DebouncedMailbox:
		receiver := make(chan Message)

		go func() {
			var lastMessage Message
			for msg := range receiver {
				lastMessage = msg
			nextMessage:
				timer := time.NewTimer(conf.Delay)
				select {
				case <-timer.C:
					to.receive <- lastMessage
				case nextMsg := <-receiver:
					lastMessage = nextMsg
					goto nextMessage
				}
			}
		}()

		deliverFunc = func(msg Message) bool {
			receiver <- msg
			return true
		}
	}

	return &mailbox{
		to:          to,
		deliverFunc: deliverFunc,
	}
}

type mailbox struct {
	to          *InPort
	deliverFunc func(msg Message) bool
}

func (m *mailbox) deliver(msg Message) bool {
	return m.deliverFunc(msg)
}
