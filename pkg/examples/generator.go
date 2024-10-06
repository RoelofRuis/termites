package examples

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"time"
)

type Generator struct {
	StringOut *termites.OutPort
	IntOut    *termites.OutPort

	sleep time.Duration
}

func NewGenerator(sleep time.Duration) *Generator {
	builder := termites.NewBuilder("Generator")

	g := &Generator{
		StringOut: termites.NewOutPort[string](builder, "String"),
		IntOut:    termites.NewOutPort[int](builder, "Int"),

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
		g.IntOut.Send(counter)
		counter++
		time.Sleep(g.sleep)
	}
}
