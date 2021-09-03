package termites

import (
	"fmt"
	"io"
)

type closeOnTeardown struct {
	closer io.Closer
}

func (c closeOnTeardown) SetEventBus(m EventBus) {
	m.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{F: c.Teardown},
	})
}

func (c closeOnTeardown) Teardown() error {
	if err := c.closer.Close(); err != nil {
		return fmt.Errorf("error closing resource: %w", err)
	}
	return nil
}