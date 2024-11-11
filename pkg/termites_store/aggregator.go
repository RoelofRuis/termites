package termites_store

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"reflect"
)

type RecordStored struct {
	RecordId RecordId `json:"record_id"`
}

type Aggregator[A any] struct {
	RecordIn *termites.InPort
	EventOut *termites.OutPort

	store Store[A]
}

func NewAggregator[A any](store Store[A]) *Aggregator[A] {
	var msg A
	dataType := reflect.TypeOf(msg)

	builder := termites.NewBuilder(fmt.Sprintf("Aggregator[%s]", dataType.Name()))

	n := &Aggregator[A]{
		RecordIn: termites.NewInPort[A](builder),
		EventOut: termites.NewOutPort[RecordStored](builder),

		store: store,
	}

	builder.OnRun(n.Run)

	return n
}

func (a *Aggregator[A]) Run(e termites.NodeControl) error {
	for msg := range a.RecordIn.Receive() {
		record := msg.Data.(A)
		id, err := a.store.Put(record)
		if err != nil {
			e.LogError("failed to store record", err)
		}
		a.EventOut.Send(RecordStored{id})
	}

	return nil
}
