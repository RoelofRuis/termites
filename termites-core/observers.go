package termites

import (
	"io"
	"log"
)

type closeOnShutdown struct {
	closer io.Closer
}

func (c closeOnShutdown) Name() string {
	return "Close On Shutdown"
}

func (c closeOnShutdown) OnNodeRegistered(Node) {}
func (c closeOnShutdown) OnGraphTeardown() {
	if err := c.closer.Close(); err != nil {
		log.Printf("Error closing: %s", err)
	}
}