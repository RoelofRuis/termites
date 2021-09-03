package termites

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Graph struct {
	name      string
	runLock   sync.Mutex
	isRunning bool

	registeredNodes map[NodeId]*node
	eventBus        *EventBus

	Close chan struct{}
}

type graphConfig struct {
	name               string
	subscribers        []EventSubscriber
	withSigtermHandler bool
	addRunner          bool
	addLogger          bool
}

func NewGraph(opts ...GraphOptions) *Graph {
	config := &graphConfig{
		name:               "",
		subscribers:        nil,
		withSigtermHandler: true,
		addRunner:          true,
		addLogger:          false,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.addLogger {
		config.subscribers = append(config.subscribers, NewConsoleLogger())
	}

	if config.addRunner {
		config.subscribers = append(config.subscribers, newRunner())
	}

	eventBus := NewEventBus()
	for _, subscriber := range config.subscribers {
		subscriber.SetEventBus(eventBus)
	}

	name := config.name
	if name == "" {
		name = "graph-" + RandomID()
	}

	g := &Graph{
		name:            name,
		registeredNodes: make(map[NodeId]*node),

		runLock:   sync.Mutex{},
		isRunning: true,
		eventBus:  eventBus,

		Close: make(chan struct{}),
	}

	if config.withSigtermHandler {
		g.setupSigtermHandler()
	}

	return g
}

func (g *Graph) ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) {
	g.Connect(out, append(opts, To(in))...)
}

func (g *Graph) Connect(out *OutPort, opts ...ConnectionOption) {
	connection, err := out.connect(opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err)) // TODO: refactor to pass panic up
	}

	g.registerNode(out.owner)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		g.registerNode(connection.mailbox.to.owner)
	}
}

func (g *Graph) Shutdown() {
	g.runLock.Lock()
	if !g.isRunning {
		return
	}
	g.isRunning = false
	fmt.Printf("Shutting down graph [%s]\n", g.name)
	g.runLock.Unlock()

	g.eventBus.Send(Event{Type: GraphTeardown})

	close(g.Close)
	fmt.Printf("Graph [%s] stopped\n", g.name)
}

func (g *Graph) setupSigtermHandler() {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		g.Shutdown()
	}()
}

func (g *Graph) registerNode(n *node) {
	_, has := g.registeredNodes[n.id]
	if !has {
		g.registeredNodes[n.id] = n

		g.eventBus.Send(Event{
			Type: NodeRegistered,
			Data: NodeRegisteredEvent{
				Node: n,
			},
		})
	}
}
