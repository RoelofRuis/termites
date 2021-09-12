package termites_dbg

import (
	"github.com/RoelofRuis/termites/termites"
)

type refReceiver struct {
	refChan        chan termites.NodeRef
	removeChan     chan termites.NodeId
	registeredRefs map[termites.NodeId]termites.NodeRef
	removedRefs    map[termites.NodeId]bool
	RefsOut        *termites.OutPort
}

func newRefReceiver() *refReceiver {
	builder := termites.NewBuilder("Ref Receiver")

	n := &refReceiver{
		refChan:        make(chan termites.NodeRef),
		removeChan:     make(chan termites.NodeId),
		registeredRefs: make(map[termites.NodeId]termites.NodeRef),
		removedRefs:    make(map[termites.NodeId]bool),
		RefsOut:        builder.OutPort("Refs", map[termites.NodeId]termites.NodeRef{}),
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
