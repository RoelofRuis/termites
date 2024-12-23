package termites_dbg

import (
	"fmt"
	"os"

	"github.com/RoelofRuis/termites/pkg/termites"
)

type Visualizer struct {
	RefsIn  *termites.InPort
	PathOut *termites.OutPort
	writer  *graphWriter
}

func NewVisualizer(fileDir string) *Visualizer {
	builder := termites.NewBuilder("Visualizer")

	n := &Visualizer{
		RefsIn:  termites.NewInPortNamed[map[termites.NodeId]termites.NodeRef](builder, "Refs"),
		PathOut: termites.NewOutPortNamed[visualizerMessage](builder, "Path"),
		writer: &graphWriter{
			rootDir:      fileDir,
			writeDotFile: false,
			version:      0,
		},
	}

	builder.OnRun(n.Run)
	builder.OnShutdown(n.Shutdown)

	return n
}

func (v *Visualizer) Run(c termites.NodeControl) error {
	if v.writer == nil {
		return fmt.Errorf("no graph writer initialized for visualizer")
	}

	for msg := range v.RefsIn.Receive() {
		refs := msg.Data.(map[termites.NodeId]termites.NodeRef)
		var nodes []termites.NodeRef
		for _, ref := range refs {
			nodes = append(nodes, ref)
		}

		path, err := v.writer.saveRoutingGraph(nodes)
		if err != nil {
			c.LogError("error creating routing graph", err)
			continue
		}

		v.PathOut.Send(visualizerMessage{path: path})
	}

	return nil
}

func (v *Visualizer) Shutdown(control termites.TeardownControl) error {
	if v.writer != nil {
		control.LogInfo(fmt.Sprintf("Cleaning up [%s]\n", v.writer.rootDir))
		return os.RemoveAll(v.writer.rootDir)
	}
	return nil
}
