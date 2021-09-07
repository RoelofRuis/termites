package termites

import (
	"fmt"
	"runtime/debug"
	"sync"
)

type NodeControl interface {
	SetSuspended()
	SetActive()
	SetError()
	LogInfo(msg string)
	LogError(msg string, err error)
}

// create via the termites.Builder
type node struct {
	name       string
	id         NodeId
	refVersion uint

	status        NodeStatus
	runningStatus NodeRunningStatus

	inPorts  []*InPort
	outPorts []*OutPort

	run      func(nodeControl NodeControl) error
	shutdown func(teardownControl TeardownControl) error

	nodeLock sync.Locker
	bus      EventSender
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

func (n *node) setStatus(s NodeStatus) {
	if n.status == s {
		return
	}

	n.nodeLock.Lock()
	n.status = s
	n.nodeLock.Unlock()
	n.sendRef()
}

func (n *node) setRunningStatus(s NodeRunningStatus) {
	if n.runningStatus == s {
		return
	}

	n.nodeLock.Lock()
	n.runningStatus = s
	n.nodeLock.Unlock()
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
	ref := n.ref()

	n.sendEvent(Event{
		Type: NodeRefUpdated,
		Data: NodeUpdatedEvent{Ref: ref},
	})
}

func (n *node) ref() NodeRef {
	n.nodeLock.Lock()
	defer n.nodeLock.Unlock()

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

	n.refVersion += 1

	return NodeRef{
		Id:            n.id,
		Version:       n.refVersion,
		Name:          n.name,
		Status:        n.status,
		RunningStatus: n.runningStatus,
		InPorts:       inPortRefs,
		OutPorts:      outPortRefs,
		RunInfo:       determineFunctionInfo(n.run),
		ShutdownInfo:  determineFunctionInfo(n.shutdown),
	}
}

func (n *node) start(bus EventBus) {
	n.nodeLock.Lock()
	if n.bus != nil {
		n.nodeLock.Unlock()
		return
	}
	n.bus = bus
	n.nodeLock.Unlock()

	if n.shutdown != nil {
		n.bus.Send(Event{
			Type: RegisterTeardown,
			Data: RegisterTeardownEvent{
				Name: n.id.String(),
				F:    n.shutdown,
			},
		})
	}

	go func() {
		n.setRunningStatus(NodeRunning)
		defer func() {
			if err := recover(); err != nil {
				n.bus.Send(LogErrorEvent(fmt.Sprintf("Node [%s] crashed: %v", n.name, err), nil))
				debug.PrintStack()
				n.SetError()
				n.setRunningStatus(NodeTerminated)
			}
		}()

		err := n.run(n)
		if err != nil {
			n.bus.Send(LogErrorEvent(fmt.Sprintf("Node [%s] exited with error", n.name), err))
			n.SetError()
			n.setRunningStatus(NodeTerminated)
			return
		}
		n.setRunningStatus(NodeTerminated)
	}()
}
