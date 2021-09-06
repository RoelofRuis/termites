package termites

import (
	"fmt"
	"io"
	"time"
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

func (c closeOnTeardown) Teardown(_ time.Duration) error {
	if err := c.closer.Close(); err != nil {
		return fmt.Errorf("error closing resource [%s] [%w]", c.name, err)
	}
	return nil
}
