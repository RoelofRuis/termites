package termites_dbg

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/RoelofRuis/termites/pkg/termites"
)

//go:embed templates/base.gohtml
var basePage string

//go:embed templates/index.gohtml
var indexPage string

//go:embed templates/nodes.gohtml
var nodesPage string

type WebController struct {
	index *template.Template
	nodes *template.Template

	uiDataLock sync.RWMutex
	uiData     UIData
	editor     CodeEditor
}

type UIData struct {
	RoutingPath string
	Nodes       []NodeInfo
}

type NodeInfo struct {
	Id          string
	Name        string
	Status      string
	Filename    string
	InPortNames []string
	Connections []ConnectionInfo
	RunInfo     termites.FunctionInfo
}

type ConnectionInfo struct {
	Id              string
	OutPortName     string
	AdapterName     string
	AdapterFilename string
	TransformInfo   termites.FunctionInfo
	InNodeName      string
	InPortName      string
}

func NewWebController() *WebController {
	return &WebController{
		index:      mustParse(basePage, indexPage),
		nodes:      mustParse(basePage, nodesPage),
		uiDataLock: sync.RWMutex{},
		uiData:     UIData{RoutingPath: "", Nodes: nil},
	}
}

func (d *WebController) SetNodes(nodes []NodeInfo) {
	d.uiDataLock.Lock()
	d.uiData.Nodes = nodes
	d.uiDataLock.Unlock()
}

func (d *WebController) SetRoutingPath(path string) {
	d.uiDataLock.Lock()
	d.uiData.RoutingPath = path
	d.uiDataLock.Unlock()
}

func (d *WebController) HandleIndex(w http.ResponseWriter, req *http.Request) {
	d.uiDataLock.RLock()
	err := d.index.ExecuteTemplate(w, "base", d.uiData)
	d.uiDataLock.RUnlock()
	if err != nil {
		panic(err)
	}
}

func (d *WebController) HandleNodes(w http.ResponseWriter, req *http.Request) {
	d.uiDataLock.RLock()
	err := d.nodes.ExecuteTemplate(w, "base", d.uiData)
	d.uiDataLock.RUnlock()
	if err != nil {
		panic(err)
	}
}

func (d *WebController) HandleOpen(w http.ResponseWriter, req *http.Request) {
	if d.editor == nil {
		http.Error(w, "no editor configured", http.StatusNotImplemented)
		return
	}

	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "no id given", http.StatusBadRequest)
		return
	}

	reses, ok := req.URL.Query()["res"]
	if !ok || len(reses[0]) < 1 {
		http.Error(w, "no res given", http.StatusBadRequest)
		return
	}

	if err := d.openResource(reses[0], ids[0]); err != nil {
		http.Error(w, "error opening resource", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/nodes", http.StatusFound)
}

func (d *WebController) openResource(resource string, id string) error {
	if d.editor == nil {
		return nil
	}

	d.uiDataLock.RLock()
	defer d.uiDataLock.RUnlock()

	if resource == "run" {
		for _, n := range d.uiData.Nodes {
			if n.Id == id && n.RunInfo.File != "" {
				if err := d.editor(n.RunInfo.File, n.RunInfo.Line); err != nil {
					return err
				}
				return nil
			}
		}
	} else if resource == "transform" {
		for _, n := range d.uiData.Nodes {
			for _, c := range n.Connections {
				if c.Id == id && c.TransformInfo.File != "" {
					if err := d.editor(c.TransformInfo.File, c.TransformInfo.Line); err != nil {
						return err
					}
					return nil
				}
			}
		}
	}

	return fmt.Errorf("resource [%s] with id [%s] not found", resource, id)
}

func mustParse(templates ...string) *template.Template {
	if len(templates) == 0 {
		panic(errors.New("at least one template must be given"))
	}
	var t = template.New("")
	var err error
	for _, data := range templates {
		t, err = t.Parse(data)
		if err != nil {
			panic(err)
		}
	}
	return t
}
