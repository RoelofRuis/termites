package termites

import (
	"fmt"
	"runtime/debug"
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

// create via the termites.Builder
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

func (n *node) setEventBus(bus EventBus) {
	n.nodeLock.Lock()
	if n.bus == nil {
		n.bus = bus
		n.bus.Send(Event{
			Type: NodeRegistered,
			Data: NodeRegisteredEvent{
				node: n,
			},
		})
		if n.shutdown != nil {
			n.bus.Send(Event{
				Type: RegisterTeardown,
				Data: RegisterTeardownEvent{
					Name: fmt.Sprintf("node-%s", Identifier(n.id).String()),
					F: n.shutdown,
				},
			})
		}
	}
	n.nodeLock.Unlock()
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

func (n *node) start() {
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
}
