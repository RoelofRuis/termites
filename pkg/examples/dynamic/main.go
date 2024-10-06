package main

import (
	"github.com/RoelofRuis/termites/pkg/examples"
	"github.com/RoelofRuis/termites/pkg/termites"
	"time"
)

func main() {
	graph := termites.NewGraph()

	generator := examples.NewGenerator(1 * time.Millisecond)
	printer := examples.NewPrinter()

	go func() {
		for {
			conn := graph.ConnectTo(generator.StringOut, printer.TextIn)
			time.Sleep(1 * time.Second)
			conn.Disconnect()
			time.Sleep(1 * time.Second)
		}
	}()

	graph.Wait()
}
