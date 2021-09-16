package termites

import (
	"fmt"
	"sync"
)

type Graph interface {
	Connect(out *OutPort, opts ...ConnectionOption) *Connection
	ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) *Connection
	Wait()
	Close()
}

type graphImpl struct {
	name     string
	eventBus *eventBus
	close    *sync.WaitGroup
}

type graphConfig struct {
	name               string
	subscribers        []EventSubscriber
	withSigtermHandler bool
	printLogs          bool
	printMessages      bool
}

func NewGraph(opts ...GraphOption) Graph {
	return newGraphImpl(opts...)
}

func newGraphImpl(opts ...GraphOption) *graphImpl {
	config := &graphConfig{
		name:               "",
		subscribers:        nil,
		withSigtermHandler: true,
		printLogs:          false,
		printMessages:      false,
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

	g := &graphImpl{
		name:     name,
		eventBus: bus,
		close:    closer,
	}

	bus.Subscribe(Exit, g.onExit)
	NewTeardownHandler(config.withSigtermHandler).SetEventBus(bus)

	if config.printLogs || config.printMessages {
		NewConsoleLogger(config.printLogs, config.printMessages).SetEventBus(bus)
	}

	for _, subscriber := range config.subscribers {
		subscriber.SetEventBus(bus)
	}

	bus.Send(LogInfo(fmt.Sprintf("Graph [%s] initialized", g.name)))
	return g
}

func (g *graphImpl) onExit(_ Event) error {
	g.close.Done()
	g.eventBus.Send(LogInfo(fmt.Sprintf("Graph [%s] closed", g.name)))
	return nil
}

func (g *graphImpl) Wait() {
	g.close.Wait()
}

func (g *graphImpl) Close() {
	g.eventBus.Send(Event{Type: Kill})
}

func (g *graphImpl) ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) *Connection {
	return g.Connect(out, append(opts, To(in))...)
}

func (g *graphImpl) Connect(out *OutPort, opts ...ConnectionOption) *Connection {
	connection, err := newConnection(out, opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	out.owner.start(g.eventBus)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		connection.mailbox.to.owner.start(g.eventBus)
	}

	return connection
}
