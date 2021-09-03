package termites

import (
	"os"
	"os/signal"
	"syscall"
)

type SigtermHandler struct{}

func NewSigtermHandler() *SigtermHandler {
	return &SigtermHandler{}
}

func (h *SigtermHandler) SetEventBus(b EventBus) {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		b.Send(Event{Type: SystemExit})
	}()
}
