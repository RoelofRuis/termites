package termites_dbg

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/RoelofRuis/termites/pkg/termites"
)

type WebUpdater struct {
	RefsIn *termites.InPort

	controller *WebController
}

func NewWebUpdater(controller *WebController) *WebUpdater {
	builder := termites.NewBuilder("Web Updater")

	n := &WebUpdater{
		RefsIn:     termites.NewInPortNamed[map[termites.NodeId]termites.NodeRef](builder, "Refs"),
		controller: controller,
	}

	builder.OnRun(n.Run)

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
							AdapterInfo:     adapterTransform,
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
		}
	}
}
