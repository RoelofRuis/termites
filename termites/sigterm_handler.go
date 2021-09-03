package termites

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type SigtermHandler struct {
	teardownFunctions map[string]func()
}

func NewSigtermHandler() *SigtermHandler {
	return &SigtermHandler{
		teardownFunctions: make(map[string]func()),
	}
}

func (h *SigtermHandler) SetEventBus(b EventBus) {
	b.Subscribe(RegisterTeardown, h.OnRegisterTeardown)

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c

		for name, f := range h.teardownFunctions {
			b.Send(LogInfoEvent(fmt.Sprintf("Running teardown for [%s]...", name)))
			f()
			b.Send(LogInfoEvent(fmt.Sprintf("Teardown for [%s] done", name)))
		}

		b.Send(Event{Type: SysExit})
	}()
}

func (h *SigtermHandler) OnRegisterTeardown(e Event) error {
	event, ok := e.Data.(RegisterTeardownEvent)
	if !ok {
		return InvalidEventError
	}
	h.teardownFunctions[event.Name] = event.F
	return nil
}
