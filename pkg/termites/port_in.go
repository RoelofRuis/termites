package termites

import (
	"reflect"
	"sync"
)

type InPort struct {
	id       InPortId
	name     string
	dataType reflect.Type
	owner    *node
	receive  chan Message

	connectionLock *sync.RWMutex
	connections    map[ConnectionId]*Connection
}

// Create via the termites.Builder
func newInPort(name string, dataType reflect.Type, owner *node) *InPort {
	return &InPort{
		id:       NewIdentifier("in-port"),
		name:     name,
		dataType: dataType,
		owner:    owner,
		receive:  make(chan Message),

		connectionLock: &sync.RWMutex{},
		connections:    make(map[ConnectionId]*Connection),
	}
}

func (p *InPort) connect(conn *Connection) {
	p.connectionLock.Lock()
	p.connections[conn.id] = conn
	p.connectionLock.Unlock()
}

func (p *InPort) disconnect(conn *Connection) {
	p.connectionLock.Lock()
	delete(p.connections, conn.id)
	p.connectionLock.Unlock()
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
