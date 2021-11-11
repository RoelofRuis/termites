package termites

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TeardownControl interface {
	GetTimeout() time.Duration
	LogInfo(msg string)
	LogError(msg string, err error)
}

type teardownControl struct {
	timeout time.Duration
	bus     EventBus
}

func (c *teardownControl) GetTimeout() time.Duration {
	return c.timeout
}

func (c *teardownControl) LogInfo(msg string) {
	c.bus.Send(LogInfo(msg))
}

func (c *teardownControl) LogError(msg string, err error) {
	c.bus.Send(LogError(msg, err))
}

type TeardownHandler struct {
	lock                 sync.Locker
	teardownFunctions    map[string]func(control TeardownControl) error
	teardownTimeout      time.Duration
	terminateOnOsSignals bool
	osChan               chan os.Signal
	killChan             chan bool
	bus                  EventBus
}

func NewTeardownHandler(terminateOnOsSignals bool) *TeardownHandler {
	return &TeardownHandler{
		lock:                 &sync.Mutex{},
		teardownFunctions:    make(map[string]func(control TeardownControl) error),
		teardownTimeout:      10 * time.Second,
		terminateOnOsSignals: terminateOnOsSignals,
		osChan:               make(chan os.Signal),
		killChan:             make(chan bool),
		bus:                  nil,
	}
}

func (h *TeardownHandler) SetEventBus(b EventBus) {
	h.bus = b
	b.Subscribe(RegisterTeardown, h.OnRegisterTeardown)
	b.Subscribe(Kill, h.OnKill)

	if h.terminateOnOsSignals {
		signal.Notify(h.osChan, os.Interrupt, syscall.SIGTERM)
	}

	go h.awaitTeardown()
}

func (h *TeardownHandler) OnKill(_ Event) error {
	h.killChan <- true
	return nil
}

func (h *TeardownHandler) OnRegisterTeardown(e Event) error {
	event, ok := e.Data.(RegisterTeardownEvent)
	if !ok {
		return InvalidEventError
	}
	h.lock.Lock()
	h.teardownFunctions[event.Name] = event.F
	h.lock.Unlock()
	return nil
}

func (h *TeardownHandler) awaitTeardown() {
	select {
	case <-h.osChan:
	case <-h.killChan:
	}

	wg := sync.WaitGroup{}
	h.lock.Lock()
	wg.Add(len(h.teardownFunctions))
	control := &teardownControl{
		timeout: 10 * time.Second,
		bus:     h.bus,
	}
	for name, f := range h.teardownFunctions {
		go func(name string, f func(TeardownControl) error) {
			h.bus.Send(LogInfo(fmt.Sprintf("Running teardown for [%s]...", name)))
			if err := f(control); err != nil {
				h.bus.Send(LogError(fmt.Sprintf("Teardown returned error"), err))
			}
			h.bus.Send(LogInfo(fmt.Sprintf("Teardown for [%s] done", name)))
			wg.Done()
		}(name, f)
	}
	h.lock.Unlock()

	await := make(chan bool)
	go func() {
		wg.Wait()
		await <- true
	}()

	select {
	case <-await:
		h.bus.Send(LogInfo("All registered teardown routines completed"))

	case <-time.NewTimer(h.teardownTimeout).C:
		h.bus.Send(LogInfo("Timeout reached for teardown routines. Forced exit"))

	}

	h.bus.Send(Event{Type: Exit})
}
