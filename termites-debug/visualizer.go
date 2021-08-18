package debug

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/RoelofRuis/termites/termites-core"
)

type Visualizer struct {
	RefsIn  *termites.InPort
	PathOut *termites.OutPort
	writer  *graphWriter
}

func NewVisualizer() *Visualizer {
	builder := termites.NewBuilder("Visualizer")

	var writer *graphWriter
	tempDir, err := ioutil.TempDir("", "vis-")
	if err == nil {
		writer = &graphWriter{
			rootDir:      tempDir,
			writeDotFile: false,
			version:      0,
		}
	}

	n := &Visualizer{
		RefsIn:  builder.InPort("Refs", map[termites.NodeId]termites.NodeRef{}),
		PathOut: builder.OutPort("Path", ""),
		writer:  writer,
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

func (v *Visualizer) Shutdown(_ time.Duration) error {
	if v.writer != nil {
		log.Printf("Cleaning up [%s]\n", v.writer.rootDir)
		return os.RemoveAll(v.writer.rootDir)
	}
	return nil
}
