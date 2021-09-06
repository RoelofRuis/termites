package termites

import (
	"errors"
	"fmt"
	"time"
)

type Connection struct {
	id      ConnectionId
	mailbox *mailbox
	adapter *Adapter
}

func (p *Connection) send(data interface{}) (error, interface{}) {
	connData := data
	if p.adapter != nil {
		var err error
		connData, err = p.adapter.transform(connData)
		if err != nil {
			return err, nil
		}
	}

	if p.mailbox == nil || connData == nil {
		return nil, connData
	}

	isDelivered := p.mailbox.deliver(Message{Data: connData})

	if !isDelivered {
		return errors.New("delivery failed"), connData
	}

	return nil, connData
}

func (p *Connection) ref() ConnectionRef {
	var adapterRef *AdapterRef = nil
	if p.adapter != nil {
		ref := p.adapter.ref()
		adapterRef = &ref
	}
	var inRef *InPortRef = nil
	if p.mailbox != nil {
		r := p.mailbox.to.ref()
		inRef = &r
	}

	return ConnectionRef{
		Id:      p.id,
		Adapter: adapterRef,
		In:      inRef,
	}
}

type connectionConfig struct {
	from    *OutPort
	to      *InPort
	adapter *Adapter
	mailbox MailboxConfig
}

func newConnection(out *OutPort, opts ...ConnectionOption) (*Connection, error) {
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

	conn := &Connection{
		id:      ConnectionId(NewIdentifier("connection")),
		mailbox: mailbox,
		adapter: config.adapter,
	}

	out.connect(conn)

	return conn, nil
}
