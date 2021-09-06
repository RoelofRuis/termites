package main

import (
	"fmt"
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_dbg"
	"time"
)

// TODO: measure resource usage/memory leaks from this example

func main() {
	graph := termites.NewGraph(
		termites_dbg.WithDebugger(4242),
	)

	generator := NewGenerator()
	printer := NewPrinter()

	go func(){
		for {
			conn := graph.ConnectTo(generator.TextOut, printer.TextIn)
			time.Sleep(1 * time.Second)
			conn.Disconnect()
			time.Sleep(1 * time.Second)
		}
	}()

	graph.Wait()
}

type Generator struct {
	TextOut *termites.OutPort
}

func NewGenerator() *Generator {
	builder := termites.NewBuilder("Generator")

	g := &Generator{
		TextOut: builder.OutPort("Generator", ""),
	}

	builder.OnRun(g.Run)

	return g
}

func (g *Generator) Run(_ termites.NodeControl) error {
	counter := 0
	for {
		g.TextOut.Send(fmt.Sprintf("%d", counter))
		counter++
		time.Sleep(1 * time.Millisecond)
	}
}

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
