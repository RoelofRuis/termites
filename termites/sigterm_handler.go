package termites

import (
	"os"
	"os/signal"
	"syscall"
)

type SigtermHandler struct {
	teardownFunctions []func() error
}

func NewSigtermHandler() *SigtermHandler {
	return &SigtermHandler{}
}

func (h *SigtermHandler) SetEventBus(b EventBus) {
	b.Subscribe(RegisterTeardown, h.OnRegisterTeardown)

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c

		for _, f := range h.teardownFunctions {
			_ = f() // TODO: not ignore error
		}

		b.Send(Event{Type: SysExit})
	}()
}

func (h *SigtermHandler) OnRegisterTeardown(e Event) error {
	event, ok := e.Data.(RegisterTeardownEvent)
	if !ok {
		return InvalidEventError
	}
	h.teardownFunctions = append(h.teardownFunctions, event.f)
	return nil
}
