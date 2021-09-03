package termites

import (
	"fmt"
	"sync"
)

type Graph struct {
	name      string
	runLock   sync.Mutex
	isRunning bool

	registeredNodes map[NodeId]*node
	eventBus        *eventBus

	Close chan struct{}
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

	g := &Graph{
		name:            name,
		registeredNodes: make(map[NodeId]*node),

		runLock:   sync.Mutex{},
		isRunning: true,
		eventBus:  bus,

		Close: make(chan struct{}),
	}

	bus.Subscribe(SysExit, g.onSysExit)

	if config.withSigtermHandler {
		g.Subscribe(NewSigtermHandler())
	}

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

func (g *Graph) onSysExit(_ Event) error {
	g.Shutdown()
	return nil
}

func (g *Graph) Shutdown() {
	g.runLock.Lock()
	if !g.isRunning {
		g.runLock.Unlock()
		return
	}
	g.isRunning = false
	g.runLock.Unlock()
	close(g.Close)
	g.eventBus.Send(LogInfoEvent(fmt.Sprintf("Graph [%s] stopped", g.name)))
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
