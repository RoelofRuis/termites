package termites

import (
	"sync"
	"time"
)

type NodeControl interface {
	SetSuspended()
	SetActive()
	SetError()
	LogInfo(msg string)
	LogError(msg string, err error)
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

	nodeLock sync.Locker
	bus      EventBus
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

func (n *node) LogInfo(msg string) {
	n.sendEvent(LogInfoEvent(msg))
}

func (n *node) LogError(msg string, err error) {
	n.sendEvent(LogErrorEvent(msg, err))
}

func (n *node) setBus(bus EventBus) {
	n.nodeLock.Lock()
	n.bus = bus
	n.nodeLock.Unlock()
}

func (n *node) setStatus(s NodeStatus) {
	if n.status == s {
		return
	}

	n.status = s
	n.sendRef()
}

func (n *node) setRunningStatus(s NodeRunningStatus) {
	if n.runningStatus == s {
		return
	}

	n.runningStatus = s
	n.sendRef()
}

func (n *node) sendEvent(e Event) {
	n.nodeLock.Lock()
	if n.bus == nil {
		n.nodeLock.Unlock()
		return
	}
	n.nodeLock.Unlock()

	n.bus.Send(e)
}

func (n *node) sendRef() {
	n.nodeLock.Lock()
	ref := n.ref()
	n.nodeLock.Unlock()

	n.sendEvent(Event{
		Type: NodeRefUpdated,
		Data: NodeUpdatedEvent{Ref: ref},
	})
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
