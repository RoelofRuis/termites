package termites

type runner struct {}

func newRunner() *runner {
	return &runner{}
}

func (r *runner) SetEventBus(b EventBus) {
	b.Subscribe(NodeRegistered, r.OnNodeRegistered)
}

func (r *runner) OnNodeRegistered(e Event) error {
	n, ok := e.Data.(NodeRegisteredEvent)
	if !ok {
		return InvalidEventError
	}

	go n.node.start()

	return nil
}
