package termites

import (
	"fmt"
	"io"
)

type closeOnShutdown struct {
	closer io.Closer
}

func (c closeOnShutdown) SetEventBus(m EventBus) {
	m.Subscribe(GraphTeardown, c.OnGraphTeardown)
}

func (c closeOnShutdown) OnGraphTeardown(_ Event) error {
	if err := c.closer.Close(); err != nil {
		return fmt.Errorf("error closing resource: %w", err)
	}
	return nil
}