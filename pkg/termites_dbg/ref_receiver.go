package termites_dbg

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"sync"
)

type refReceiver struct {
	RefsOut *termites.OutPort

	lock           sync.Mutex
	registeredRefs map[termites.NodeId]termites.NodeRef
	removedRefs    map[termites.NodeId]bool
}

func newRefReceiver() *refReceiver {
	builder := termites.NewBuilder("Ref Receiver")

	n := &refReceiver{
		RefsOut:        termites.NewOutPortNamed[map[termites.NodeId]termites.NodeRef](builder, "Refs"),
		registeredRefs: make(map[termites.NodeId]termites.NodeRef),
		removedRefs:    make(map[termites.NodeId]bool),
	}

	return n
}

func (r *refReceiver) onNodeRefUpdated(e termites.Event) error {
	msg, ok := e.Data.(termites.NodeUpdatedEvent)
	if !ok {
		return termites.InvalidEventError
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	current, has := r.registeredRefs[msg.Ref.Id]
	if has && msg.Ref.Version < current.Version {
		return nil
	}
	if _, has := r.removedRefs[msg.Ref.Id]; has {
		return nil
	}
	r.registeredRefs[msg.Ref.Id] = msg.Ref

	refsToSend := make(map[termites.NodeId]termites.NodeRef, len(r.registeredRefs))
	for id, r := range r.registeredRefs {
		refsToSend[id] = r
	}

	r.RefsOut.Send(refsToSend)

	return nil
}

func (r *refReceiver) onNodeStopped(e termites.Event) error {
	msg, ok := e.Data.(termites.NodeStoppedEvent)
	if !ok {
		return termites.InvalidEventError
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.removedRefs[msg.Id] = true
	delete(r.registeredRefs, msg.Id)
	refsToSend := make(map[termites.NodeId]termites.NodeRef, len(r.registeredRefs))
	for id, r := range r.registeredRefs {
		refsToSend[id] = r
	}

	r.RefsOut.Send(refsToSend)

	return nil
}
