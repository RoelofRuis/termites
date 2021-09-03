package termites

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type runner struct {
	shutdownFuncs   []func(timeout time.Duration) error
	shutdownTimeout time.Duration
	bus             EventBus
}

func newRunner() *runner {
	return &runner{
		shutdownFuncs:   nil,
		shutdownTimeout: time.Second * 10,
	}
}

func (r *runner) SetEventBus(b EventBus) {
	b.Subscribe(NodeRegistered, r.OnNodeRegistered)
	b.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{f: r.Teardown},
	})
	r.bus = b
}

func (r *runner) OnNodeRegistered(e Event) error {
	n, ok := e.Data.(NodeRegisteredEvent)
	if !ok {
		return InvalidEventError
	}

	if n.node.shutdown != nil {
		r.shutdownFuncs = append(r.shutdownFuncs, n.node.shutdown)
	}

	go func(node *node) {
		node.setBus(r.bus)
		node.setRunningStatus(NodeRunning)
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Node [%s] crashed\nPanic: %s\n", node.name, err)
				debug.PrintStack()
				node.SetError()
				node.setRunningStatus(NodeTerminated)
			}
		}()

		err := node.run(node)
		if err != nil {
			fmt.Printf("Node [%s] exited with error\nError was: %v\n", node.name, err)
			node.SetError()
			node.setRunningStatus(NodeTerminated)
			return
		}
		node.setRunningStatus(NodeTerminated)
	}(n.node)

	return nil
}

func (r *runner) Teardown() error {
	// TODO: logging etc
	fmt.Printf("Stopping Runner...\n")

	wg := sync.WaitGroup{}
	wg.Add(len(r.shutdownFuncs))

	for _, shutdown := range r.shutdownFuncs {
		go func(f func(timeout time.Duration) error) {
			err := f(r.shutdownTimeout)
			if err != nil {
				fmt.Printf("Error when shutting down node: %v", err.Error())
			}
			wg.Done()
		}(shutdown)
	}

	await := make(chan bool)
	go func() {
		wg.Wait()
		await <- true
	}()

	select {
	case <-await:
		fmt.Printf("All shutdown routines completed\n")

	case <-time.NewTimer(r.shutdownTimeout).C:
		fmt.Printf("Shutdown timeout reached, force exit\n")
	}

	return nil
}
