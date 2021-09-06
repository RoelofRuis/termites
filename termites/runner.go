package termites

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type runner struct {
	shutdownRegistry *shutdownRegistry
	bus              EventBus
}

func newRunner() *runner {
	return &runner{
		shutdownRegistry: newShutdownRegistry(10 * time.Second),
		bus:              nil,
	}
}

func (r *runner) SetEventBus(b EventBus) {
	r.shutdownRegistry.bus = b
	r.bus = b
	r.bus.Subscribe(NodeRegistered, r.OnNodeRegistered)
	r.bus.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{Name: "runner", F: r.shutdownRegistry.shutdownAll},
	})
}

func (r *runner) OnNodeRegistered(e Event) error {
	n, ok := e.Data.(NodeRegisteredEvent)
	if !ok {
		return InvalidEventError
	}

	r.shutdownRegistry.register(n.node)

	go run(r.bus, n.node)

	return nil
}

func run(bus EventBus, node *node) {
	node.setRunningStatus(NodeRunning)
	defer func() {
		if err := recover(); err != nil {
			bus.Send(LogErrorEvent(fmt.Sprintf("Node [%s] crashed: %v", node.name, err), nil))
			debug.PrintStack()
			node.SetError()
			node.setRunningStatus(NodeTerminated)
		}
	}()

	err := node.run(node)
	if err != nil {
		bus.Send(LogErrorEvent(fmt.Sprintf("Node [%s] exited with error", node.name), err))
		node.SetError()
		node.setRunningStatus(NodeTerminated)
		return
	}
	node.setRunningStatus(NodeTerminated)
}

type shutdownRegistry struct {
	lock            sync.Locker
	shutdownFuncs   []func(timeout time.Duration) error
	shutdownTimeout time.Duration
	bus             EventBus
}

func newShutdownRegistry(shutdownTimeout time.Duration) *shutdownRegistry {
	return &shutdownRegistry{
		lock:            &sync.Mutex{},
		shutdownFuncs:   nil,
		shutdownTimeout: shutdownTimeout,
		bus:             nil,
	}
}

func (r *shutdownRegistry) register(n *node) {
	if n.shutdown != nil {
		r.lock.Lock()
		r.shutdownFuncs = append(r.shutdownFuncs, n.shutdown)
		r.lock.Unlock()
	}
}

func (r *shutdownRegistry) shutdownAll() {
	wg := sync.WaitGroup{}
	r.lock.Lock()
	wg.Add(len(r.shutdownFuncs))

	for _, shutdown := range r.shutdownFuncs {
		go func(f func(timeout time.Duration) error) {
			err := f(r.shutdownTimeout)
			if err != nil {
				r.bus.Send(LogErrorEvent("Error when shutting down node", err))
			}
			wg.Done()
		}(shutdown)
	}
	r.lock.Unlock()

	await := make(chan bool)
	go func() {
		wg.Wait()
		await <- true
	}()

	select {
	case <-await:
		r.bus.Send(LogInfoEvent("All registered node shutdown routines completed"))

	case <-time.NewTimer(r.shutdownTimeout).C:
		r.bus.Send(LogInfoEvent("Shutdown timeout reached for node shutdown routines. Forced exit."))
	}
}
