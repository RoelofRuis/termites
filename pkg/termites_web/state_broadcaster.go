package termites_web

import (
	"encoding/json"
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

// AsMutationFor ensures that a Mutation message for State S can be passed to a node with that receiving type.
func AsMutationFor[S State]() termites.ConnectionOption {
	return termites.Via(func(a any) (Mutation[S], error) {
		cast, ok := a.(Mutation[S])
		if !ok {
			return nil, fmt.Errorf("value is not a Mutation[S]")
		}
		return cast, nil
	})
}

type Mutation[S State] interface {
	Mutate(state S) error
}

type State interface {
	Read() (json.RawMessage, error)
}

type StateBroadcaster[S State] struct {
	ConnectionIn *termites.InPort
	MutationsIn  *termites.InPort
	MessageOut   *termites.OutPort

	state     S
	stateData json.RawMessage
}

func NewStateBroadcaster[S State](initialState S) *StateBroadcaster[S] {
	builder := termites.NewBuilder("StateBroadcaster")

	t := &StateBroadcaster[S]{
		ConnectionIn: termites.NewInPort[ClientConnection](builder),
		MutationsIn:  termites.NewInPort[Mutation[S]](builder),
		MessageOut:   termites.NewOutPort[ClientMessage](builder),

		state: initialState,
	}

	builder.OnRun(t.Run)

	return t
}

func (v *StateBroadcaster[S]) Run(c termites.NodeControl) error {
	var err error
	v.stateData, err = v.state.Read()
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-v.ConnectionIn.Receive():
			connection := msg.Data.(ClientConnection)

			clientMessage, err := NewClientMessageFor("state/full", connection.Id, v.stateData)
			if err != nil {
				c.LogError("failed to marshal state", err)
				continue
			}

			v.MessageOut.Send(clientMessage)

		case msg := <-v.MutationsIn.Receive():
			mutation := msg.Data.(Mutation[S])

			oldState := v.stateData

			err = mutation.Mutate(v.state)
			if err != nil {
				c.LogError("failed to apply mutation", err)
				continue
			}

			v.stateData, err = v.state.Read()
			if err != nil {
				c.LogError("failed to retrieve state", err)
				continue
			}

			patch, err := jsonpatch.CreateMergePatch(oldState, v.stateData)
			if err != nil {
				c.LogError("failed to create merge patch", err)
				continue
			}

			clientMessage, _ := NewClientMessage("state/patch", patch)

			v.MessageOut.Send(clientMessage)
		}
	}
}
