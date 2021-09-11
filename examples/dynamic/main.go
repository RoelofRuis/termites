package main

import (
	"github.com/RoelofRuis/termites/examples"
	"github.com/RoelofRuis/termites/termites"
	"time"
)

// TODO: measure resource usage/memory leaks from this example

func main() {
	graph := termites.NewGraph()

	generator := examples.NewGenerator(1 * time.Millisecond)
	printer := examples.NewPrinter()

	go func() {
		for {
			conn := graph.ConnectTo(generator.TextOut, printer.TextIn)
			time.Sleep(1 * time.Second)
			conn.Disconnect()
			time.Sleep(1 * time.Second)
		}
	}()

	graph.Wait()
}
