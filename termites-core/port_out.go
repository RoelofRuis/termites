package termites

import (
	"sync"
)

type OutPort struct {
	id       OutPortId
	name     string
	dataType string
	owner    *node

	connections []connection
}

// Create via the termites.Builder
func newOutPort(name string, dataType string, owner *node) *OutPort {
	return &OutPort{
		id:          OutPortId(NewIdentifier("out-port")),
		name:        name,
		dataType:    dataType,
		owner:       owner,
		connections: nil,
	}
}

func (p *OutPort) Send(data interface{}) {
	if len(p.connections) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for _, conn := range p.connections {
		wg.Add(1)
		go func(conn connection) {
			err, sentData := conn.send(data)
			if p.owner.bus != nil {
				// TODO: clean up and push as much as possible down into event
				toName := ""
				toPortName := ""
				if conn.mailbox != nil {
					toName = conn.mailbox.to.owner.name
					toPortName = conn.mailbox.to.name
				}
				adapterName := ""
				if conn.adapter != nil {
					adapterName = conn.adapter.name
				}
				p.owner.bus.Send(Event{
					Type: MessageSent,
					Data: MessageSentEvent{
						FromName:     p.owner.name,
						FromPortName: p.name,
						ToName:       toName,
						ToPortName:   toPortName,
						AdapterName:  adapterName,
						Data:         sentData,
						Error:        err,
					},
				})
			}
			wg.Done()
		}(conn)
	}
	wg.Wait()
}

func (p *OutPort) ref() OutPortRef {
	var connections []ConnectionRef
	for _, conn := range p.connections {
		connections = append(connections, conn.ref())
	}
	return OutPortRef{
		Id:          p.id,
		Name:        p.name,
		Connections: connections,
	}
}
