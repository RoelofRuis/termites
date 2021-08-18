package debug

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/RoelofRuis/termites/termites-core"
)

//go:embed templates/layout.gohtml
var layoutPage string

//go:embed templates/index.gohtml
var indexPage string

//go:embed templates/nodes.gohtml
var nodesPage string

type WebUI struct {
	HttpPort  int
	StaticDir string
	DataLock  sync.Mutex
	UIData    UIData
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
	RunInfo     *termites.FunctionInfo
}

type ConnectionInfo struct {
	Id              string
	OutPortName     string
	AdapterName     string
	AdapterFilename string
	TransformInfo   *termites.FunctionInfo
	InNodeName      string
	InPortName      string
}

func NewWebUI(httpPort int, staticDir string) *WebUI {
	return &WebUI{
		HttpPort:  httpPort,
		StaticDir: staticDir,
		DataLock:  sync.Mutex{},
		UIData: UIData{
			RoutingPath: "",
			Nodes:       nil,
		},
	}
}

func (d *WebUI) run() {
	if d.HttpPort == 0 {
		return
	}

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var err error
		t := template.New("template")
		t, err = t.Parse(layoutPage)
		if err != nil {
			panic(err)
		}
		t, err = t.Parse(indexPage)
		if err != nil {
			panic(err)
		}
		d.DataLock.Lock()
		err = t.ExecuteTemplate(w, "layout", d.UIData)
		d.DataLock.Unlock()
		if err != nil {
			panic(err)
		}
	})

	m.HandleFunc("/nodes", func(w http.ResponseWriter, req *http.Request) {
		t := template.New("template")
		t, err := t.Parse(layoutPage)
		if err != nil {
			panic(err)
		}
		t, err = t.Parse(nodesPage)
		if err != nil {
			panic(err)
		}
		d.DataLock.Lock()
		err = t.ExecuteTemplate(w, "layout", d.UIData)
		d.DataLock.Unlock()
		if err != nil {
			panic(err)
		}
	})

	m.HandleFunc("/open", func(w http.ResponseWriter, req *http.Request) {
		ids, ok := req.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reses, ok := req.URL.Query()["res"]
		if !ok || len(reses[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := d.openResource(reses[0], ids[0]); err != nil {
			log.Printf("error: %s", err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Redirect(w, req, "/nodes", http.StatusFound)
	})

	fs := http.FileServer(http.Dir(d.StaticDir))
	m.Handle("/static/", http.StripPrefix("/static", fs))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", d.HttpPort), m); err != nil {
		panic(err)
	}
}

func (d *WebUI) openResource(resource string, id string) error {
	d.DataLock.Lock()
	defer d.DataLock.Unlock()

	if resource == "run" {
		for _, n := range d.UIData.Nodes {
			if n.Id == id {
				if err := Open(n.RunInfo.File, n.RunInfo.Line); err != nil {
					return err
				}
				return nil
			}
		}

	} else if resource == "transform" {
		for _, n := range d.UIData.Nodes {
			for _, c := range n.Connections {
				if c.Id == id {
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
