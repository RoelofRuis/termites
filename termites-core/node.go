package termites

import (
	"time"
)

type Node interface {
	SetNodeRefChannel(chan<- NodeRef)
	SetMessageRefChannel(chan<- MessageRef)
	getNode() *node
}

type NodeControl interface { // TODO: add error logging via control
	SetSuspended()
	SetActive()
	SetError()
}

type node struct {
	name string
	id   NodeId

	status        NodeStatus
	runningStatus NodeRunningStatus

	inPorts  []*InPort
	outPorts []*OutPort

	run      func(nodeController NodeControl) error
	shutdown func(timeout time.Duration) error

	nodeRefChannel    chan<- NodeRef
	messageRefChannel chan<- MessageRef
}

func (n *node) SetNodeRefChannel(c chan<- NodeRef) {
	n.nodeRefChannel = c
	n.updateRef()
}

func (n *node) SetMessageRefChannel(c chan<- MessageRef) {
	n.messageRefChannel = c
}

func (n *node) getNode() *node {
	return n
}

func (n *node) SetSuspended() {
	n.setStatus(NodeSuspended)
}

func (n *node) SetActive() {
	n.setStatus(NodeActive)
}

func (n *node) SetError() {
	n.setStatus(NodeError)
}

func (n *node) setStatus(s NodeStatus) {
	if n.status == s {
		return
	}

	n.status = s
	n.updateRef()
}

func (n *node) setRunningStatus(s NodeRunningStatus) {
	if n.runningStatus == s {
		return
	}

	n.runningStatus = s
	n.updateRef()
}

func (n *node) updateRef() {
	if n.nodeRefChannel == nil {
		return
	}

	n.nodeRefChannel <- n.ref()
}

func (n *node) ref() NodeRef {
	inPortRefs := make(map[InPortId]InPortRef)
	for _, in := range n.inPorts {
		inPortRefs[in.id] = InPortRef{
			Id:   in.id,
			Name: in.name,
		}
	}

	outPortRefs := make(map[OutPortId]OutPortRef)
	for _, out := range n.outPorts {
		outPortRefs[out.id] = out.ref()
	}

	return NodeRef{
		Id:            n.id,
		Name:          n.name,
		Status:        n.status,
		RunningStatus: n.runningStatus,
		InPorts:       inPortRefs,
		OutPorts:      outPortRefs,
		RunInfo:       determineFunctionInfo(n.run),
		ShutdownInfo:  determineFunctionInfo(n.shutdown),
	}
}
