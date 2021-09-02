package termites

import (
	"log"
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
			if p.owner.messageRefChannel != nil {
				// TODO: probably clean up this logic and push as much as possible down to logger
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
				ref := MessageRef{
					fromName:     p.owner.name,
					fromPortName: p.name,
					toName:       toName,
					toPortName:   toPortName,
					adapterName:  adapterName,
					data:         sentData,
					error:        err,
				}
				select {
				case p.owner.messageRefChannel <- ref:
				default:
					log.Print("DROPPED MESSAGE REF\n")
				}
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
