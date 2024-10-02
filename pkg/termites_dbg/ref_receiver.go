package termites_dbg

import (
	"github.com/RoelofRuis/termites/pkg/termites"
)

type refReceiver struct {
	RefsOut *termites.OutPort

	refChan        chan termites.NodeRef
	removeChan     chan termites.NodeId
	registeredRefs map[termites.NodeId]termites.NodeRef
	removedRefs    map[termites.NodeId]bool
}

func newRefReceiver() *refReceiver {
	builder := termites.NewBuilder("Ref Receiver")

	n := &refReceiver{
		RefsOut:        termites.NewOutPort[map[termites.NodeId]termites.NodeRef](builder, "Refs"),
		refChan:        make(chan termites.NodeRef),
		removeChan:     make(chan termites.NodeId),
		registeredRefs: make(map[termites.NodeId]termites.NodeRef),
		removedRefs:    make(map[termites.NodeId]bool),
	}

	builder.OnRun(n.run)

	return n
}

func (r *refReceiver) run(_ termites.NodeControl) error {
	for {
		select {
		case ref := <-r.refChan:
			current, has := r.registeredRefs[ref.Id]
			if has && ref.Version < current.Version {
				continue
			}
			if _, has = r.removedRefs[ref.Id]; has {
				continue
			}

			r.registeredRefs[ref.Id] = ref

			r.sendAll()

		case id := <-r.removeChan:
			r.removedRefs[id] = true
			delete(r.registeredRefs, id)
			r.sendAll()
		}
	}
}

func (r *refReceiver) sendAll() {
	refsToSend := make(map[termites.NodeId]termites.NodeRef, len(r.registeredRefs))
	for id, r := range r.registeredRefs {
		refsToSend[id] = r
	}

	r.RefsOut.Send(refsToSend)
}
