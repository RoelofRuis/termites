package termites

import (
	"sync"
)

type OutPort struct {
	id       OutPortId
	name     string
	dataType string
	owner    *node

	connectionLock *sync.RWMutex
	connections    map[ConnectionId]*Connection
}

// Create via the termites.Builder
func newOutPort(name string, dataType string, owner *node) *OutPort {
	return &OutPort{
		id:             OutPortId(NewIdentifier("out-port")),
		name:           name,
		dataType:       dataType,
		owner:          owner,
		connectionLock: &sync.RWMutex{},
		connections:    make(map[ConnectionId]*Connection),
	}
}

func (p *OutPort) Send(data interface{}) {
	p.connectionLock.RLock()
	defer p.connectionLock.RUnlock()

	if len(p.connections) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for _, conn := range p.connections {
		wg.Add(1)
		go func(conn *Connection) {
			err, _ := conn.send(data)
			p.sendMessageEvent(conn, err)
			wg.Done()
		}(conn)
	}
	wg.Wait()
}

func (p *OutPort) connect(conn *Connection) {
	p.connectionLock.Lock()
	p.connections[conn.id] = conn
	p.connectionLock.Unlock()
	p.owner.sendRef()
}

func (p *OutPort) disconnect(conn *Connection) {
	p.connectionLock.Lock()
	delete(p.connections, conn.id)
	p.connectionLock.Unlock()
}

func (p *OutPort) ref() OutPortRef {
	var connections []ConnectionRef
	p.connectionLock.RLock()
	for _, conn := range p.connections {
		connections = append(connections, conn.ref())
	}
	p.connectionLock.RUnlock()
	return OutPortRef{
		Id:          p.id,
		Name:        p.name,
		Connections: connections,
	}
}

func (p *OutPort) sendMessageEvent(conn *Connection, err error) {
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
	p.owner.sendEvent(Event{
		Type: MessageSent,
		Data: MessageSentEvent{
			FromName:     p.owner.name,
			FromPortName: p.name,
			ToName:       toName,
			ToPortName:   toPortName,
			AdapterName:  adapterName,
			Error:        err,
		},
	})
}
