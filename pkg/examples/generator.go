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
		StringOut: termites.NewOutPort[string](builder, "String"),
		CountOut:  termites.NewOutPort[Count](builder, "Int"),

		sleep: sleep,
	}

	builder.OnRun(g.Run)

	return g
}

func (g *Generator) Run(_ termites.NodeControl) error {
	counter := 0
	for {
		text := fmt.Sprintf("%d", counter)
		g.StringOut.Send(text)
		g.CountOut.Send(Count{counter})
		counter++
		time.Sleep(g.sleep)
	}
}
