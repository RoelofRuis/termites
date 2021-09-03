package termites

import (
	"fmt"
	"time"
)

type connectionFactory struct {}

type connectionConfig struct {
	from    *OutPort
	to      *InPort
	adapter *Adapter
	mailbox MailboxConfig
}

func newConnectionFactory() *connectionFactory {
	return &connectionFactory{}
}

func (f *connectionFactory) newConnection(out *OutPort, opts ...ConnectionOption) (*connection, error) {
	if out == nil {
		return nil, fmt.Errorf("cannot connect nil out port")
	}

	config := &connectionConfig{from: out, to: nil, adapter: nil, mailbox: nil}

	for _, opt := range opts {
		opt(config)
	}

	if config.from == nil {
		return nil, fmt.Errorf("no out port configured")
	}

	if config.adapter == nil && config.to == nil {
		return nil, fmt.Errorf("no adapter and in port configured, at least one should be given")
	}

	if config.adapter == nil && config.from.dataType != config.to.dataType {
		return nil, fmt.Errorf("out port [%s:%s (%s)] and in port [%s (%s)] have differing data types\n",
			config.from.owner.name,
			config.from.name,
			config.from.dataType,
			config.to.name,
			config.to.dataType,
		)
	}

	if config.adapter != nil && config.from.dataType != config.adapter.inDataType {
		return nil, fmt.Errorf("out port [%s:%s (%s)] and adapter [%s (%s)] have differing data types\n",
			config.from.owner.name,
			config.from.name,
			config.from.dataType,
			config.adapter.name,
			config.adapter.inDataType,
		)
	}

	if config.adapter != nil && config.to == nil && config.adapter.outDataType != determineDataType(nil) {
		return nil, fmt.Errorf("adapter [%s (%s)] is not connected to in and must have 'nil' data out\n",
			config.adapter.name,
			config.adapter.outDataType,
		)
	}

	if config.adapter != nil && config.to != nil && config.adapter.outDataType != config.to.dataType {
		return nil, fmt.Errorf(
			"adapter [%s (%s)] and in port [%s:%s (%s)] have differing data types\n",
			config.adapter.name,
			config.adapter.outDataType,
			config.to.owner.name,
			config.to.name,
			config.to.dataType,
		)
	}

	var mailbox *mailbox = nil
	if config.to != nil {
		if config.mailbox == nil {
			config.mailbox = &NormalMailbox{ReceiveTimeout: 1 * time.Second}
		}
		mailbox = mailboxFromConfig(config.to, config.mailbox)
	}

	return &connection{
		id:      ConnectionId(NewIdentifier("connection")),
		mailbox: mailbox,
		adapter: config.adapter,
	}, nil
}
