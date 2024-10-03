package termites_dbg

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"path/filepath"

	"github.com/RoelofRuis/termites/pkg/termites"
)

var VisualizerAdapter = termites.NewAdapter(
	"Visualizer Path",
	func(visualizerPath string) (termites.JsonPartialData, error) {
		_, filename := filepath.Split(visualizerPath)

		msg, err := json.Marshal(struct {
			Path string `json:"path"`
		}{
			Path: filepath.Join("/dbg-static/", filename),
		})
		if err != nil {
			return termites.JsonPartialData{}, err
		}
		return termites.JsonPartialData{Key: "routing_graph", Data: msg}, nil
	},
)

var ClientAdapter = termites.NewAdapter(
	"Client Adapter",
	func(bytes []byte) (termites_web.ClientMessage, error) {
		return termites_web.ClientMessage{Data: bytes}, nil
	},
)
