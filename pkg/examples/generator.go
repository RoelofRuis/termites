package examples

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"time"
)

// Count serves as a general state struct that can be sent throughout the program.
type Count struct {
	Count int `json:"count"`
}

type Generator struct {
	StringOut *termites.OutPort
	CountOut  *termites.OutPort

	sleep time.Duration
}

func NewGenerator(sleep time.Duration) *Generator {
	builder := termites.NewBuilder("Generator")

	g := &Generator{
		StringOut: termites.NewOutPort[string](builder),
		CountOut:  termites.NewOutPort[Count](builder),

		sleep: sleep,
	}

	builder.OnRun(g.Run)

	return g
}

func (g *Generator) Run(e termites.NodeControl) error {
	e.LogInfo("Starting the generator...")

	counter := 0
	for {
		e.LogInfo("Generating the next number...")
		text := fmt.Sprintf("%d", counter)
		g.StringOut.Send(text)
		g.CountOut.Send(Count{counter})
		counter++
		time.Sleep(g.sleep)
	}
}
