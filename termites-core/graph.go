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
	name        string
	runLock     sync.Mutex
	blockingRun bool
	isRunning   bool

	hooks             []GraphHook
	connectionFactory *connectionFactory

	nodeRegistry NodeRegistry

	shutdown chan struct{}
}

type GraphHook interface {
	Setup(registry NodeRegistry)
	Teardown()
}

type graphConfig struct {
	hooks              []GraphHook
	withSigtermHandler bool
	addRunner          bool
	addLogger          bool
	blockingRun        bool
}

func NewGraph(opts ...GraphOptions) *Graph {
	config := &graphConfig{
		hooks:              nil,
		withSigtermHandler: true,
		addRunner:          true,
		addLogger:          false,
		blockingRun:        true,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.addLogger {
		config.hooks = append(config.hooks, newLogger())
	}

	if config.addRunner {
		config.hooks = append(config.hooks, newRunner())
	}

	g := &Graph{
		name:              "graph-" + uuid.New().String(),
		hooks:             config.hooks,
		connectionFactory: newConnectionFactory(),

		nodeRegistry: newNodeRegistry(),
		runLock:      sync.Mutex{},
		blockingRun:  config.blockingRun,
		isRunning:    false,

		shutdown: nil,
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
	if g.isRunning {
		panic(fmt.Errorf("graph is already running"))
	}

	connection, err := g.connectionFactory.newConnection(out, opts...)
	if err != nil {
		panic(fmt.Errorf("node connection error: %w", err))
	}

	out.connections = append(out.connections, *connection)
	g.nodeRegistry.register(out.owner)
	if connection.mailbox != nil && connection.mailbox.to != nil {
		g.nodeRegistry.register(connection.mailbox.to.owner)
	}
}

func (g *Graph) Run() {
	g.runLock.Lock()
	if g.isRunning {
		panic(fmt.Errorf("cannot run graph twice"))
	}
	g.isRunning = true
	g.runLock.Unlock()

	for _, h := range g.hooks {
		h.Setup(g.nodeRegistry)
	}
	g.shutdown = make(chan struct{})
	if g.blockingRun {
		<-g.shutdown
	}
}

func (g *Graph) Shutdown() {
	g.runLock.Lock()
	defer g.runLock.Unlock()
	if !g.isRunning {
		return
	}
	for _, h := range g.hooks {
		h.Teardown()
	}
	close(g.shutdown)
	g.isRunning = false
}

func (g *Graph) setupSigtermHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		g.Shutdown()
	}()
}
