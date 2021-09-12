package examples

import (
	"fmt"
	"github.com/RoelofRuis/termites/termites"
)

type Printer struct {
	TextIn *termites.InPort
}

func NewPrinter() *Printer {
	builder := termites.NewBuilder("Printer")

	p := &Printer{
		TextIn: builder.InPort("Text", ""),
	}

	builder.OnRun(p.Run)

	return p
}

func (p *Printer) Run(_ termites.NodeControl) error {
	for msg := range p.TextIn.Receive() {
		fmt.Printf("PRINT: %s\n", msg.Data)
	}
	return nil
}
