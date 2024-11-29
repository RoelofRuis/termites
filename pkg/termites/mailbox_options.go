package termites

import (
	"errors"
	"time"
)

type mailboxConfig struct {
	// capacity sets the internal mailbox buffer capacity so messages are stored temporarily if they cannot be
	// immediately handled.
	capacity int

	// receiveTimeout ensures message delivery within the given timeout. Otherwise, messages are dropped.
	// If > 0  : use that value as the timeout duration.
	// If == 0 : do not use a timeout, but block until the message is accepted.
	// If < 0  : non-blocking, dropping a message immediately if it cannot be delivered.
	receiveTimeout time.Duration

	// debounceDelay sets a delay time for a debounce mechanism that holds the most recently delivered message and posts
	// it if no new messages have been received within the delay.
	debounceDelay time.Duration

	// errorWhenDropped is the error when a message is dropped. Nil should be used to silently discard dropped messages.
	errorWhenDropped error
}

type MailboxOption func(conf *mailboxConfig)

func WithTimeout(timeout time.Duration) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.receiveTimeout = timeout
		conf.errorWhenDropped = errors.New("delivery timed out")
	}
}

func WithBuffer(capacity int) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.capacity = capacity
	}
}

func WithDebounce(debounceDelay time.Duration) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.receiveTimeout = 0
		conf.debounceDelay = debounceDelay
		conf.errorWhenDropped = nil
	}
}

func WithDiscard() MailboxOption {
	return func(conf *mailboxConfig) {
		conf.receiveTimeout = -1
		conf.errorWhenDropped = errors.New("message discarded")
	}
}

func WithSilentDiscard() MailboxOption {
	return func(conf *mailboxConfig) {
		conf.receiveTimeout = -1
		conf.errorWhenDropped = nil
	}
}
