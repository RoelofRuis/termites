package examples

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"time"
)

type Generator struct {
	TextOut  *termites.OutPort
	BytesOut *termites.OutPort

	sleep time.Duration
}

func NewGenerator(sleep time.Duration) *Generator {
	builder := termites.NewBuilder("Generator")

	g := &Generator{
		TextOut:  builder.OutPort("Text", ""),
		BytesOut: builder.OutPort("Bytes", []byte{}),

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
		g.BytesOut.Send([]byte(text))
		counter++
		time.Sleep(g.sleep)
	}
}
