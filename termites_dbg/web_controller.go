package termites_dbg

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/RoelofRuis/termites/termites"
)

//go:embed templates/layout.gohtml
var layoutPage string

//go:embed templates/index.gohtml
var indexPage string

//go:embed templates/nodes.gohtml
var nodesPage string

type WebController struct {
	index *template.Template
	nodes *template.Template

	uiDataLock sync.RWMutex
	uiData     UIData
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
		index:      mustParse(layoutPage, indexPage),
		nodes:      mustParse(layoutPage, nodesPage),
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
	err := d.index.ExecuteTemplate(w, "layout", d.uiData)
	d.uiDataLock.RUnlock()
	if err != nil {
		panic(err)
	}
}

func (d *WebController) HandleNodes(w http.ResponseWriter, req *http.Request) {
	d.uiDataLock.RLock()
	err := d.nodes.ExecuteTemplate(w, "layout", d.uiData)
	d.uiDataLock.RUnlock()
	if err != nil {
		panic(err)
	}
}

func (d *WebController) HandleOpen(w http.ResponseWriter, req *http.Request) {
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
		log.Printf("error: %s", err.Error())
	}

	http.Redirect(w, req, "/nodes", http.StatusFound)
}

func (d *WebController) openResource(resource string, id string) error {
	d.uiDataLock.RLock()
	defer d.uiDataLock.RUnlock()

	if resource == "run" {
		for _, n := range d.uiData.Nodes {
			if n.Id == id && n.RunInfo.File != "" {
				if err := Open(n.RunInfo.File, n.RunInfo.Line); err != nil {
					return err
				}
				return nil
			}
		}
	} else if resource == "transform" {
		for _, n := range d.uiData.Nodes {
			for _, c := range n.Connections {
				if c.Id == id && c.TransformInfo.File != "" {
					if err := Open(c.TransformInfo.File, c.TransformInfo.Line); err != nil {
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
