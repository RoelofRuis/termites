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

func (c closeOnTeardown) Teardown(control TeardownControl) error {
	if err := c.closer.Close(); err != nil {
		control.LogError(fmt.Sprintf("error closing resource [%s]", c.name), err)
		return err
	}
	return nil
}
