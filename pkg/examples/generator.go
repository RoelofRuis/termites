package examples

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"time"
)

type Generator struct {
	TextOut   *termites.OutPort
	ClientOut *termites.OutPort

	sleep time.Duration
}

func NewGenerator(sleep time.Duration) *Generator {
	builder := termites.NewBuilder("Generator")

	g := &Generator{
		TextOut:   termites.NewOutPort[string](builder, "Text"),
		ClientOut: termites.NewOutPort[termites_web.ClientMessage](builder, "Bytes"),

		sleep: sleep,
	}

	builder.OnRun(g.Run)

	return g
}

func (g *Generator) Run(_ termites.NodeControl) error {
	counter := 0
	for {
		text := fmt.Sprintf("%d", counter)
		g.TextOut.Send(text)
		g.ClientOut.Send(termites_web.ClientMessage{Data: []byte(text)})
		counter++
		time.Sleep(g.sleep)
	}
}
