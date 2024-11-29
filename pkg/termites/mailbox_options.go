package termites

import "time"

type mailboxConfig struct {
	capacity       int
	receiveTimeout time.Duration
	debounceDelay  time.Duration
	dropOnOverflow bool
}

type MailboxOption func(conf *mailboxConfig)

func WithReceiveTimeout(timeout time.Duration) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.receiveTimeout = timeout
	}
}

func WithCapacity(capacity int) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.capacity = capacity
	}
}

func WithDebounceDelay(debounceDelay time.Duration) MailboxOption {
	return func(conf *mailboxConfig) {
		conf.debounceDelay = debounceDelay
	}
}

func WithDropOnOverflow() MailboxOption {
	return func(conf *mailboxConfig) {
		conf.dropOnOverflow = true
	}
}
