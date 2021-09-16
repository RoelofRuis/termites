package termites_dbg

import (
	"fmt"
	"os"

	"github.com/RoelofRuis/termites/termites"
)

type Visualizer struct {
	RefsIn  *termites.InPort
	PathOut *termites.OutPort
	writer  *graphWriter
}

func NewVisualizer(fileDir string) *Visualizer {
	builder := termites.NewBuilder("Visualizer")

	n := &Visualizer{
		RefsIn:  builder.InPort("Refs", map[termites.NodeId]termites.NodeRef{}),
		PathOut: builder.OutPort("Path", ""),
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

func (v *Visualizer) Run(_ termites.NodeControl) error {
	if v.writer == nil {
		return fmt.Errorf("no graph writer initialized for visualizer")
	}

	for msg := range v.RefsIn.Receive() {
		refs := msg.Data.(map[termites.NodeId]termites.NodeRef)
		var nodes []termites.NodeRef
		for _, ref := range refs {
			nodes = append(nodes, ref)
		}

		path := v.writer.saveRoutingGraph(nodes)

		if path != "" {
			v.PathOut.Send(path)
		}
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
