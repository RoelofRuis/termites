package termites_state

import (
	"github.com/RoelofRuis/termites/pkg/termites"
)

type Action[S any] interface {
	Apply(s *S)
}

type Mutator[S any] struct {
	ActionIn *termites.InPort
	StateOut *termites.OutPort

	state *S
}

func NewMutator[S any](state *S) *Mutator[S] {
	builder := termites.NewBuilder("Mutator")

	node := &Mutator[S]{
		ActionIn: termites.NewInPort[Action[S]](builder, "Action In"),
		StateOut: termites.NewOutPort[S](builder, "State Out"),

		state: state,
	}

	builder.OnRun(node.Run)

	return node
}

func (m *Mutator[S]) Run(_ termites.NodeControl) error {
	m.StateOut.Send(*m.state)

	for msg := range m.ActionIn.Receive() {
		action := msg.Data.(Action[S])
		action.Apply(m.state)
		m.StateOut.Send(*m.state)
	}

	return nil
}
