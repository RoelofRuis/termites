package termites

import (
	"fmt"
	"time"
)

type ConnectionOption func(conn *connectionConfig)

func Via(adapter *Adapter) ConnectionOption {
	if adapter == nil {
		panic(fmt.Errorf("invalid connection option: adapter cannot be nil"))
	}
	return func(conn *connectionConfig) {
		conn.adapter = adapter
	}
}

func To(in *InPort) ConnectionOption {
	if in == nil {
		panic(fmt.Errorf("invalid connection option: in port cannot be nil"))
	}
	return func(conn *connectionConfig) {
		conn.to = in
	}
}

func WithMailbox(conf MailboxConfig) ConnectionOption {
	return func(conn *connectionConfig) {
		conn.mailbox = conf
	}
}

func WithSmallCapacityMailbox() ConnectionOption {
	return WithMailbox(&CapacityMailbox{Capacity: 10, ReceiveTimeout: 1 * time.Second})
}
