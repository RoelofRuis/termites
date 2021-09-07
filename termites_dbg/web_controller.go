package termites_dbg

import (
	"fmt"
	"github.com/RoelofRuis/termites/termites"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type WebController struct {
	PathIn *termites.InPort
	RefsIn *termites.InPort

	ui *WebUI
}

func NewWebController(httpPort int) *WebController {
	staticDir, err := ioutil.TempDir("", "web-ui-")
	if err != nil {
		panic(err)
	}

	builder := termites.NewBuilder("Web Controller")

	n := &WebController{
		PathIn: builder.InPort("Visualizer Path", ""),
		RefsIn: builder.InPort("Refs", map[termites.NodeId]termites.NodeRef{}),
		ui:     NewWebUI(httpPort, staticDir),
	}

	builder.OnRun(n.Run)
	builder.OnShutdown(n.Shutdown)

	return n
}

func (d *WebController) Run(_ termites.NodeControl) error {
	go d.ui.run()

	for {
		select {
		case msg := <-d.RefsIn.Receive():
			refs := msg.Data.(map[termites.NodeId]termites.NodeRef)
			var nodes []NodeInfo
			for _, ref := range refs {
				_, file := path.Split(ref.RunInfo.File)

				status := "?"
				if ref.Status == termites.NodeActive {
					status = "active"
				} else if ref.Status == termites.NodeSuspended {
					status = "suspended"
				} else if ref.Status == termites.NodeError {
					status = "error"
				}

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
						var adapterTranform *termites.FunctionInfo
						var adapterFileName string
						if c.Adapter != nil {
							adapterName = c.Adapter.Name
							adapterTranform = c.Adapter.TransformInfo
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
							TransformInfo:   adapterTranform,
							InNodeName:      inNodeName,
							InPortName:      inPortName,
						})
					}
				}

				nodes = append(nodes, NodeInfo{
					Id:          fmt.Sprintf("%x", ref.Id),
					Name:        ref.Name,
					Status:      status,
					Filename:    file,
					InPortNames: inPortNames,
					Connections: connections,
					RunInfo:     ref.RunInfo,
				})
			}
			d.ui.DataLock.Lock()
			sort.SliceStable(nodes, func(i, j int) bool {
				return strings.Compare(nodes[i].Name, nodes[j].Name) < 0
			})
			d.ui.UIData.Nodes = nodes
			d.ui.DataLock.Unlock()

		case msg := <-d.PathIn.Receive():
			visualizerPath := msg.Data.(string)
			src, err := os.Open(visualizerPath)
			if err != nil {
				fmt.Printf("Error opening source: %s", err.Error())
				continue
			}
			_, filename := filepath.Split(visualizerPath)
			staticPath := filepath.Join(d.ui.StaticDir, filename)
			dst, err := os.OpenFile(staticPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				fmt.Printf("Error opening dst: %s", err.Error())
				continue
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				fmt.Printf("Error copying routing: %s", err.Error())
			}
			d.ui.DataLock.Lock()
			d.ui.UIData.RoutingPath = filepath.Join("/static", filename)
			d.ui.DataLock.Unlock()
		}
	}
}

func (d *WebController) Shutdown(control termites.TeardownControl) error {
	control.LogInfo(fmt.Sprintf("Cleaning up [%s]\n", d.ui.StaticDir))
	return os.RemoveAll(d.ui.StaticDir)
}
