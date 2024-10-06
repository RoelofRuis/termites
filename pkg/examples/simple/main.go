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

	graph.ConnectTo(generator.StringOut, printer.TextIn)

	graph.Wait()
}
