package termites

import (
	"fmt"
	"sync"
)

type Graph interface {
	Connect(out *OutPort, opts ...ConnectionOption)
	ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption)
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
	addRunner          bool
	addConsoleLogger   bool
}

func NewGraph(opts ...GraphOptions) Graph {
	return newGraphImpl(opts...)
}

func newGraphImpl(opts ...GraphOptions) *graphImpl {
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

	g := &graphImpl{
		name:     name,
		eventBus: bus,
		close:    closer,
	}

	bus.Subscribe(Exit, g.onExit)
	NewTeardownHandler(config.withSigtermHandler).SetEventBus(bus)

	if config.addConsoleLogger {
		NewConsoleLogger().SetEventBus(bus)
	}

	if config.addRunner {
		newRunner().SetEventBus(bus)
	}

	for _, subscriber := range config.subscribers {
		subscriber.SetEventBus(bus)
	}

	return g
}

func (g *graphImpl) onExit(_ Event) error {
	g.close.Done()
	g.eventBus.Send(LogInfoEvent(fmt.Sprintf("Graph [%s] closed", g.name)))
	return nil
}

func (g *graphImpl) Wait() {
	g.close.Wait()
}

func (g *graphImpl) Close() {
	g.eventBus.Send(Event{Type: Kill})
}

func (g *graphImpl) ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) {
	g.Connect(out, append(opts, To(in))...)
}

func (g *graphImpl) Connect(out *OutPort, opts ...ConnectionOption) {
	connection, err := newConnection(out, opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	out.owner.setBus(g.eventBus)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		connection.mailbox.to.owner.setBus(g.eventBus)
	}
}
