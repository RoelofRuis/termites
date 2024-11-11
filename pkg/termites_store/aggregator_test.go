package termites_store

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	graph := termites.NewGraph()

	store := NewInMemory[int]()

	input := termites.NewInspectableNode[int]("input")
	aggregator := NewAggregator[int](store)
	output := termites.NewInspectableNode[RecordStored]("output")

	graph.Connect(input.Out, aggregator.RecordIn)
	graph.Connect(aggregator.EventOut, output.In)

	values := []int{1, 1, 2, 3, 5, 8}

	for idx, value := range values {
		input.Send <- value

		evt, err := output.ReceiveWithin(time.Second)
		if err != nil {
			t.Fatal(err)
		}
		if string(evt.RecordId) != fmt.Sprintf("%d", idx) {
			t.Errorf("expected id 0, got %s", evt.RecordId)
		}
	}

	stored := store.GetAll()
	if len(stored) != len(values) {
		t.Errorf("expected %d records, got %d", len(values), len(stored))
	}
	for idx, value := range values {
		if stored[idx] != value {
			t.Errorf("expected value %d, got %d", value, stored[idx])
		}
	}
}
