package termites

import (
	"fmt"
	"sync"
)

// A Graph manages a collection nodes and their connections and serves as the builder of the connections.
type Graph struct {
	name        string
	eventBus    *eventBus
	close       *sync.WaitGroup
	connections []*Connection
}

func NewGraph(opts ...GraphOption) *Graph {
	config := &graphConfig{
		name:               "",
		subscribers:        nil,
		withSigtermHandler: true,
		printLogs:          false,
		printMessages:      false,
		deferredStart:      false,
	}

	for _, opt := range opts {
		opt(config)
	}

	name := config.name
	if name == "" {
		name = "graph-" + RandomID()
	}

	bus := newEventBus()
	closer := &sync.WaitGroup{}
	closer.Add(1)

	var connections []*Connection = nil
	if config.deferredStart {
		connections = []*Connection{}
	}

	g := &Graph{
		name:        name,
		eventBus:    bus,
		close:       closer,
		connections: connections,
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

func (g *Graph) onExit(_ Event) error {
	g.close.Done()
	g.eventBus.Send(LogInfo(fmt.Sprintf("Graph [%s] closed", g.name)))
	return nil
}

func (g *Graph) Wait() {
	for _, conn := range g.connections {
		g.start(conn)
	}
	g.connections = nil

	g.close.Wait()
}

func (g *Graph) Close() {
	g.eventBus.Send(Event{Type: Kill})
}

func (g *Graph) Connect(out *OutPort, in *InPort, opts ...ConnectionOption) *Connection {
	return g.ConnectBy(out, append(opts, To(in))...)
}

func (g *Graph) ConnectBy(out *OutPort, opts ...ConnectionOption) *Connection {
	connection, err := newConnection(out, opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	if g.connections == nil {
		g.start(connection)
	} else {
		g.connections = append(g.connections, connection)
	}

	return connection
}

func (g *Graph) start(c *Connection) {
	c.from.owner.start(g.eventBus)
	if c.mailbox != nil && c.mailbox.to != nil {
		c.mailbox.to.owner.start(g.eventBus)
	}
}
