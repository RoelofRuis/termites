package termites

import (
	"fmt"
	"sync"
)

type Graph struct {
	name            string
	registeredNodes map[NodeId]*node
	eventBus        *eventBus
	close           *sync.WaitGroup
}

type graphConfig struct {
	name               string
	subscribers        []EventSubscriber
	withSigtermHandler bool
	addRunner          bool
	addConsoleLogger   bool
}

func NewGraph(opts ...GraphOptions) *Graph {
	config := &graphConfig{
		name:               "",
		subscribers:        nil,
		withSigtermHandler: true,
		addRunner:          true,
		addConsoleLogger:   false,
	}

	for _, opt := range opts {
		opt(config)
	}

	name := config.name
	if name == "" {
		name = "graph-" + RandomID()
	}

	bus := NewEventBus()
	closer := &sync.WaitGroup{}
	closer.Add(1)

	g := &Graph{
		name:            name,
		registeredNodes: make(map[NodeId]*node),

		eventBus: bus,

		close: closer,
	}

	bus.Subscribe(Exit, g.onExit)
	g.Subscribe(NewTeardownHandler(config.withSigtermHandler))

	if config.addConsoleLogger {
		g.Subscribe(NewConsoleLogger())
	}

	if config.addRunner {
		g.Subscribe(newRunner())
	}

	for _, subscriber := range config.subscribers {
		g.Subscribe(subscriber)
	}

	return g
}

func (g *Graph) onExit(_ Event) error {
	g.close.Done()
	g.eventBus.Send(LogInfoEvent(fmt.Sprintf("Graph [%s] closed", g.name)))
	return nil
}

func (g *Graph) Wait() {
	g.close.Wait()
}

func (g *Graph) Kill() {
	g.eventBus.Send(Event{Type: Kill})
}

func (g *Graph) Subscribe(sub EventSubscriber) {
	sub.SetEventBus(g.eventBus)
}

func (g *Graph) ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) {
	g.Connect(out, append(opts, To(in))...)
}

func (g *Graph) Connect(out *OutPort, opts ...ConnectionOption) {
	connection, err := out.connect(opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	g.registerNode(out.owner)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		g.registerNode(connection.mailbox.to.owner)
	}
}

func (g *Graph) registerNode(n *node) {
	_, has := g.registeredNodes[n.id]
	if !has {
		g.registeredNodes[n.id] = n

		g.eventBus.Send(Event{
			Type: NodeRegistered,
			Data: NodeRegisteredEvent{
				node: n,
			},
		})
	}
}
