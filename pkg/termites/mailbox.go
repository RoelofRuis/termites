package termites

import (
	"errors"
	"time"
)

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
	var deliverFunc func(msg Message) error
	switch conf := c.(type) {
	case *NormalMailbox:
		deliverFunc = func(msg Message) error {
			ticker := time.NewTimer(conf.ReceiveTimeout)
			select {
			case <-ticker.C:
				return errors.New("delivery timed out")

			case to.receive <- msg:
				return nil
			}
		}

	case *CapacityMailbox:
		messages := make(chan Message, conf.Capacity)

		go func() {
			for msg := range messages {
				to.receive <- msg
			}
		}()

		deliverFunc = func(msg Message) error {
			ticker := time.NewTimer(conf.ReceiveTimeout)
			select {
			case <-ticker.C:
				return errors.New("delivery timed out")

			case messages <- msg:
				return nil
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

		deliverFunc = func(msg Message) error {
			receiver <- msg
			return nil
		}
	}

	return &mailbox{
		to:      to,
		deliver: deliverFunc,
	}
}

type mailbox struct {
	to      *InPort
	deliver func(msg Message) error
}
