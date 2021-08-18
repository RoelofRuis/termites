package termites

import (
	"errors"
)

type connection struct {
	id      ConnectionId
	mailbox *mailbox
	adapter *Adapter
}

func (p *connection) send(data interface{}) (error, interface{}) {
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

func (p *connection) ref() ConnectionRef {
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
