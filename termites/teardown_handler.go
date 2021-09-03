package termites

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type TeardownHandler struct {
	teardownFunctions    map[string]func()
	terminateOnOsSignals bool
	osChan               chan os.Signal
	killChan             chan bool
}

func NewTeardownHandler(terminateOnOsSignals bool) *TeardownHandler {
	return &TeardownHandler{
		teardownFunctions:    make(map[string]func()),
		terminateOnOsSignals: terminateOnOsSignals,
		osChan:               make(chan os.Signal),
		killChan:             make(chan bool),
	}
}

func (h *TeardownHandler) SetEventBus(b EventBus) {
	b.Subscribe(RegisterTeardown, h.OnRegisterTeardown)
	b.Subscribe(Kill, h.OnKill)

	if h.terminateOnOsSignals {
		signal.Notify(h.osChan, os.Interrupt, syscall.SIGTERM)
	}

	go func() {
		select {
		case <-h.osChan:
		case <-h.killChan:
		}

		for name, f := range h.teardownFunctions {
			b.Send(LogInfoEvent(fmt.Sprintf("Running teardown for [%s]...", name)))
			f()
			b.Send(LogInfoEvent(fmt.Sprintf("Teardown for [%s] done", name)))
		}

		b.Send(Event{Type: Exit})
	}()
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
	h.teardownFunctions[event.Name] = event.F
	return nil
}
