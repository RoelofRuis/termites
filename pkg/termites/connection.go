package termites

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Message struct {
	Data interface{}
}

type Connection struct {
	id      ConnectionId
	from    *OutPort
	mailbox *mailbox
	adapter *adapter
}

func (c *Connection) Disconnect() {
	c.from.disconnect(c)
}

func (c *Connection) send(data interface{}) {
	connData := data
	if c.adapter != nil {
		var err error
		connData, err = c.adapter.transform(connData)
		if errors.Is(err, SkipElement) {
			return
		}
		if err != nil {
			c.notifySent(err)
			return
		}
	}

	if c.mailbox == nil {
		return
	}

	if err := c.mailbox.deliver(Message{Data: connData}); err != nil {
		c.notifySent(err)
	}

	c.notifySent(nil)
}

// notifySent notify the graph that a message was sent.
func (c *Connection) notifySent(err error) {
	toName := ""
	toPortName := ""
	if c.mailbox != nil {
		toName = c.mailbox.to.owner.name
		toPortName = c.mailbox.to.name
	}
	adapterName := ""
	if c.adapter != nil {
		adapterName = c.adapter.name
	}
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	c.from.owner.sendEvent(Event{
		Type: MessageSent,
		Data: MessageSentEvent{
			FromName:     c.from.owner.name,
			FromPortName: c.from.name,
			ToName:       toName,
			ToPortName:   toPortName,
			AdapterName:  adapterName,
			Error:        errorMsg,
		},
	})
}

func (c *Connection) ref() ConnectionRef {
	var adapterRef *AdapterRef = nil
	if c.adapter != nil {
		ref := c.adapter.ref()
		adapterRef = &ref
	}
	var inRef *InPortRef = nil
	if c.mailbox != nil {
		r := c.mailbox.to.ref()
		inRef = &r
	}

	return ConnectionRef{
		Id:      c.id,
		Adapter: adapterRef,
		In:      inRef,
	}
}

func newConnection(out *OutPort, opts ...ConnectionOption) (*Connection, error) {
	if out == nil {
		return nil, fmt.Errorf("cannot connect nil out port")
	}

	config := &connectionConfig{from: out, to: nil, adapter: nil, mailboxOptions: nil}

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
		return nil, fmt.Errorf("out port [%s:%s (%s)] and in port [%s:%s (%s)] have differing data types\n",
			config.from.owner.name,
			config.from.name,
			config.from.dataType,
			config.to.owner.name,
			config.to.name,
			config.to.dataType,
		)
	}

	if config.adapter != nil && config.adapter.inDataType != reflect.TypeFor[interface{}]() && config.from.dataType != config.adapter.inDataType {
		return nil, fmt.Errorf("out port [%s:%s (%s)] and adapter [%s (%s)] have differing data types\n",
			config.from.owner.name,
			config.from.name,
			config.from.dataType,
			config.adapter.name,
			config.adapter.inDataType,
		)
	}

	if config.adapter != nil && config.to == nil && config.adapter.outDataType != reflect.TypeOf(nil) {
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

	var mbox *mailbox = nil
	if config.to != nil {
		if config.mailboxOptions == nil {
			config.mailboxOptions = []MailboxOption{WithTimeout(1 * time.Second)} // TODO: is this a sensible default?
		}
		mbox = newMailbox(config.to, config.mailboxOptions...)
	}

	conn := &Connection{
		id:      NewIdentifier("connection"),
		from:    out,
		mailbox: mbox,
		adapter: config.adapter,
	}

	if mbox != nil {
		mbox.to.connect(conn)
	}

	out.connect(conn)

	return conn, nil
}
