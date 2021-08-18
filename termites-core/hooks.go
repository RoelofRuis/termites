package termites

import (
	"io"
	"log"
)

type closeOnShutdown struct {
	closer io.Closer
}

func (c closeOnShutdown) Setup(_ NodeRegistry) {}
func (c closeOnShutdown) Teardown() {
	if err := c.closer.Close(); err != nil {
		log.Printf("Error closing: %s", err)
	}
}
