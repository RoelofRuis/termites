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

func (b *Builder) InPort(name string, exampleMsg interface{}) *InPort {
	dataType := reflect.TypeOf(exampleMsg)

	in := newInPort(name, dataType, b.node)
	b.node.inPorts = append(b.node.inPorts, in)

	return in
}

func (b *Builder) OutPort(name string, exampleMsg interface{}) *OutPort {
	dataType := reflect.TypeOf(exampleMsg)

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
