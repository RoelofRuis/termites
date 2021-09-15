package termites_dbg

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/RoelofRuis/termites/termites"
)

type WebUpdater struct {
	PathIn *termites.InPort
	RefsIn *termites.InPort

	controller *WebController
	staticDir  string
}

func NewWebUpdater(staticDir string, controller *WebController) *WebUpdater {
	builder := termites.NewBuilder("Web Updater")

	n := &WebUpdater{
		PathIn:     builder.InPort("Visualizer Path", ""),
		RefsIn:     builder.InPort("Refs", map[termites.NodeId]termites.NodeRef{}),
		controller: controller,
		staticDir:  staticDir,
	}

	builder.OnRun(n.Run)
	builder.OnShutdown(n.Shutdown)

	return n
}

func (d *WebUpdater) Run(_ termites.NodeControl) error {
	for {
		select {
		case msg := <-d.RefsIn.Receive():
			refs := msg.Data.(map[termites.NodeId]termites.NodeRef)
			var nodes []NodeInfo
			for _, ref := range refs {
				_, file := path.Split(ref.RunInfo.File)

				var inPortNames []string
				for _, i := range ref.InPorts {
					inPortNames = append(inPortNames, i.Name)
				}

				var outPortNames []string
				var connections []ConnectionInfo
				for _, i := range ref.OutPorts {
					outPortNames = append(outPortNames, i.Name)
					for _, c := range i.Connections {
						var adapterName string
						var adapterTransform termites.FunctionInfo
						var adapterFileName string
						if c.Adapter != nil {
							adapterName = c.Adapter.Name
							adapterTransform = c.Adapter.TransformInfo
							_, adapterFileName = path.Split(c.Adapter.TransformInfo.File)
						}
						var inPortName, inNodeName string
						if c.In != nil {
							inPortName = c.In.Name
							for _, r := range refs {
								for _, inPort := range r.InPorts {
									if inPort.Id == c.In.Id {
										inNodeName = r.Name
										break
									}
								}
							}
						}

						connections = append(connections, ConnectionInfo{
							Id:              fmt.Sprintf("%x", c.Id),
							OutPortName:     i.Name,
							AdapterName:     adapterName,
							AdapterFilename: adapterFileName,
							TransformInfo:   adapterTransform,
							InNodeName:      inNodeName,
							InPortName:      inPortName,
						})
					}
				}

				nodes = append(nodes, NodeInfo{
					Id:          fmt.Sprintf("%x", ref.Id),
					Name:        ref.Name,
					Status:      "active",
					Filename:    file,
					InPortNames: inPortNames,
					Connections: connections,
					RunInfo:     ref.RunInfo,
				})
			}
			sort.SliceStable(nodes, func(i, j int) bool {
				return strings.Compare(nodes[i].Name, nodes[j].Name) < 0
			})
			d.controller.SetNodes(nodes)

		case msg := <-d.PathIn.Receive():
			visualizerPath := msg.Data.(string)
			src, err := os.Open(visualizerPath)
			if err != nil {
				fmt.Printf("Error opening source: %s", err.Error())
				continue
			}
			_, filename := filepath.Split(visualizerPath)
			staticPath := filepath.Join(d.staticDir, filename)
			dst, err := os.OpenFile(staticPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				fmt.Printf("Error opening dst: %s", err.Error())
				continue
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				fmt.Printf("Error copying routing: %s", err.Error())
			}
			d.controller.SetRoutingPath(filepath.Join("/termites_dbg-static/", filename))
		}
	}
}

func (d *WebUpdater) Shutdown(control termites.TeardownControl) error {
	control.LogInfo(fmt.Sprintf("Cleaning up [%s]\n", d.staticDir))
	return os.RemoveAll(d.staticDir)
}