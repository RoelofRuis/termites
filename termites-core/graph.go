package termites

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

type Graph struct {
	name      string
	runLock   sync.Mutex
	isRunning bool

	registeredNodes   map[NodeId]*node
	observers         []GraphObserver
	connectionFactory *connectionFactory

	Close chan struct{}
}

type GraphObserver interface {
	Name() string
	OnNodeRegistered(n Node)
	OnGraphTeardown()
}

type graphConfig struct {
	name               string
	observers          []GraphObserver
	withSigtermHandler bool
	addRunner          bool
	addLogger          bool
}

func NewGraph(opts ...GraphOptions) *Graph {
	config := &graphConfig{
		name:               "",
		observers:          nil,
		withSigtermHandler: true,
		addRunner:          true,
		addLogger:          false,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.addLogger {
		config.observers = append(config.observers, newLogger())
	}

	if config.addRunner {
		config.observers = append(config.observers, newRunner())
	}

	name := config.name
	if name == "" {
		name = "graph-" + uuid.New().String()
	}

	g := &Graph{
		name:              name,
		registeredNodes:   make(map[NodeId]*node),
		observers:         config.observers,
		connectionFactory: newConnectionFactory(),

		runLock:   sync.Mutex{},
		isRunning: true,

		Close: make(chan struct{}),
	}

	if config.withSigtermHandler {
		g.setupSigtermHandler()
	}

	for _, o := range g.observers {
		fmt.Printf("Graph [%s] has observer [%s]\n", g.name, o.Name())
	}

	return g
}

func (g *Graph) ConnectTo(out *OutPort, in *InPort, opts ...ConnectionOption) {
	g.Connect(out, append(opts, To(in))...)
}

func (g *Graph) Connect(out *OutPort, opts ...ConnectionOption) {
	connection, err := g.connectionFactory.newConnection(out, opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	out.connections = append(out.connections, *connection)

	g.registerNode(out.owner)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		g.registerNode(connection.mailbox.to.owner)
	}

	// update ref to ensure listeners are notified of new connections
	// TODO: rethink this, if we split it up so that connections are their own concept, we might not have to do this here
	out.owner.updateRef()
}

func (g *Graph) Shutdown() {
	g.runLock.Lock()
	if !g.isRunning {
		return
	}
	g.isRunning = false
	fmt.Printf("Shutting down graph [%s]\n", g.name)
	g.runLock.Unlock()

	for _, o := range g.observers {
		o.OnGraphTeardown()
	}

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

		for _, o := range g.observers {
			o.OnNodeRegistered(n)
		}
	}
}
