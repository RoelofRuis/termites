package termites_dbg

import (
	"github.com/RoelofRuis/termites/termites"
)

type refReceiver struct {
	refChan        chan termites.NodeRef
	registeredRefs map[termites.NodeId]termites.NodeRef
	RefsOut        *termites.OutPort
}

func newRefReceiver() *refReceiver {
	builder := termites.NewBuilder("Ref Receiver")

	n := &refReceiver{
		refChan:        make(chan termites.NodeRef),
		registeredRefs: make(map[termites.NodeId]termites.NodeRef),
		RefsOut:        builder.OutPort("Refs", map[termites.NodeId]termites.NodeRef{}),
	}

	builder.OnRun(n.run)

	return n
}

func (r *refReceiver) run(_ termites.NodeControl) error {
	for ref := range r.refChan {
		current, has := r.registeredRefs[ref.Id]
		if has && ref.Version < current.Version {
			continue
		}
		r.registeredRefs[ref.Id] = ref

		refsToSend := make(map[termites.NodeId]termites.NodeRef, len(r.registeredRefs))
		for id, r := range r.registeredRefs {
			refsToSend[id] = r
		}

		r.RefsOut.Send(refsToSend)
	}

	return nil
}
