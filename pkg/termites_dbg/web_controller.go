package termites_dbg

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/RoelofRuis/termites/pkg/termites"
)

//go:embed templates/base.gohtml
var basePage string

//go:embed templates/index.gohtml
var indexPage string

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
	AdapterInfo     termites.FunctionInfo
	InNodeName      string
	InPortName      string
}

func NewWebController() *WebController {
	return &WebController{
		index:      mustParse(basePage, indexPage),
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

	parts := strings.Split(ids[0], ":")
	if len(parts) != 2 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if parts[0] != "runner" && parts[0] != "adapter" {
		http.Error(w, "invalid resource type given", http.StatusBadRequest)
		return
	}

	if err := d.openResource(parts[0], parts[1]); err != nil {
		http.Error(w, "error opening resource", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func (d *WebController) openResource(resource string, id string) error {
	if d.editor == nil {
		return nil
	}

	d.uiDataLock.RLock()
	defer d.uiDataLock.RUnlock()

	if resource == "runner" {
		for _, n := range d.uiData.Nodes {
			if n.Id == id {
				if err := open(n.RunInfo, d.editor); err != nil {
					return err
				}
				return nil
			}
		}
	} else if resource == "adapter" {
		for _, n := range d.uiData.Nodes {
			for _, c := range n.Connections {
				if c.Id == id {
					if err := open(c.AdapterInfo, d.editor); err != nil {
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
