package termites

import (
	"fmt"
	"io"
)

type closeOnTeardown struct {
	name   string
	closer io.Closer
	bus    EventBus
}

func (c closeOnTeardown) SetEventBus(b EventBus) {
	b.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{Name: c.name, F: c.Teardown},
	})
	c.bus = b
}

func (c closeOnTeardown) Teardown() {
	if err := c.closer.Close(); err != nil {
		c.bus.Send(LogErrorEvent(fmt.Sprintf("error closing resource [%s]", c.name), err))
	}
}
