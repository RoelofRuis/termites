package termites

import (
	"reflect"
	"sync"
)

type Builder struct {
	node *node
}

func NewBuilder(name string) Builder {
	return Builder{
		node: &node{
			id:         NewIdentifier("node"),
			name:       name,
			refVersion: 0,

			inPorts:  nil,
			outPorts: nil,

			run:      nil,
			shutdown: nil,

			nodeLock: &sync.Mutex{},
			bus:      nil,
		},
	}
}

func NewInPort[A any](b Builder) *InPort {
	var msg A
	dataType := reflect.TypeOf(msg)

	if dataType == nil {
		dataType = reflect.TypeFor[A]()
	}

	in := newInPort(dataType.Name(), dataType, b.node)
	b.node.inPorts = append(b.node.inPorts, in)

	return in
}

func NewOutPort[A any](b Builder) *OutPort {
	var msg A
	dataType := reflect.TypeOf(msg)

	if dataType == nil {
		dataType = reflect.TypeFor[A]()
	}

	out := newOutPort(dataType.Name(), dataType, b.node)
	b.node.outPorts = append(b.node.outPorts, out)

	return out
}

func NewInPortNamed[A any](b Builder, name string) *InPort {
	var msg A
	dataType := reflect.TypeOf(msg)

	in := newInPort(name, dataType, b.node)
	b.node.inPorts = append(b.node.inPorts, in)

	return in
}

func NewOutPortNamed[A any](b Builder, name string) *OutPort {
	var msg A
	dataType := reflect.TypeOf(msg)

	out := newOutPort(name, dataType, b.node)
	b.node.outPorts = append(b.node.outPorts, out)

	return out
}

func (b *Builder) OnRun(f func(control NodeControl) error) {
	b.node.run = f
}

func (b *Builder) OnShutdown(f func(control TeardownControl) error) {
	b.node.shutdown = f
}
