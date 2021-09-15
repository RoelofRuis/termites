package main

import (
	examples2 "github.com/RoelofRuis/termites/examples"
	"github.com/RoelofRuis/termites/termites"
	"time"
)

func main() {
	graph := termites.NewGraph()

	generator := examples2.NewGenerator(1 * time.Millisecond)
	printer := examples2.NewPrinter()

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
