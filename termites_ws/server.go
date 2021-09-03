package termites_ws

import (
	"embed"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type WebsocketServer struct {
	router   *mux.Router
	registry ClientRegistry
	address  string
	rootDir  string
}

func NewWebsocketServer(
	registry ClientRegistry,
	address string,
	rootDir string,
) *WebsocketServer {
	return &WebsocketServer{
		router:   mux.NewRouter(),
		registry: registry,
		address:  address,
		rootDir:  rootDir,
	}
}

func (s *WebsocketServer) Run() {
	connector := NewWebsocketConnector(s.registry)

	s.router.Path("/").Methods("GET").HandlerFunc(s.handleIndex)
	s.router.Path("/ws").Methods("GET").HandlerFunc(connector.ConnectWebsocket)

	fileServer := http.FileServer(http.Dir(s.rootDir))
	s.router.PathPrefix("/static/").Methods("GET").Handler(http.StripPrefix("/static/", fileServer))

	embeddedServer := http.FileServer(http.FS(embeddedJS))
	s.router.PathPrefix("/embedded/").Methods("GET").Handler(http.StripPrefix("/embedded/", embeddedServer))

	go func() {
		err := http.ListenAndServe(s.address, s.router)
		if err != nil {
			log.Println(err)
		}
	}()
}

//go:embed connect.js
var embeddedJS embed.FS

func (s *WebsocketServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(s.rootDir, "index.html"))
}
