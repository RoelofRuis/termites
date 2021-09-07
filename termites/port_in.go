package termites

type InPort struct {
	id       InPortId
	name     string
	dataType string
	owner    *node
	receive  chan Message
}

// Create via the termites.Builder
func newInPort(name string, dataType string, owner *node) *InPort {
	return &InPort{
		id:       NewIdentifier("in-port"),
		name:     name,
		dataType: dataType,
		owner:    owner,
		receive:  make(chan Message),
	}
}

func (p *InPort) Receive() <-chan Message {
	return p.receive
}

func (p *InPort) ref() InPortRef {
	return InPortRef{
		Id:   p.id,
		Name: p.name,
	}
}
