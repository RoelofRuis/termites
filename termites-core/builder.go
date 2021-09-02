package termites

import (
	"time"
)

type Builder struct {
	node *node
}

func NewBuilder(name string) Builder {
	return Builder{
		node: &node{
			id:            NodeId(NewIdentifier("node")),
			name:          name,
			status:        NodeActive,
			runningStatus: NodePreStarted,
			inPorts:       nil,
			outPorts:      nil,
			run:           nil,
			shutdown:      nil,
		},
	}
}

func (b *Builder) InPort(name string, exampleMsg interface{}) *InPort {
	dataType := determineDataType(exampleMsg)

	in := newInPort(name, dataType, b.node)
	b.node.inPorts = append(b.node.inPorts, in)

	return in
}

func (b *Builder) OutPort(name string, exampleMsg interface{}) *OutPort {
	dataType := determineDataType(exampleMsg)

	out := newOutPort(name, dataType, b.node)
	b.node.outPorts = append(b.node.outPorts, out)

	return out
}

func (b *Builder) OnRun(f func(nodeController NodeControl) error) {
	b.node.run = f
}

func (b *Builder) OnShutdown(f func(timeout time.Duration) error) {
	b.node.shutdown = f
}
