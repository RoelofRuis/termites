package termites_dbg

import (
	"encoding/json"
	"path/filepath"

	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_web"
)

var VisualizerAdapter = termites.NewAdapter(
	"Visualizer Path",
	func(visualizerPath string) (termites_web.JsonPartialData, error) {
		_, filename := filepath.Split(visualizerPath)

		msg, err := json.Marshal(struct {
			Path string `json:"path"`
		}{
			Path: filepath.Join("/dbg-static/", filename),
		})
		if err != nil {
			return termites_web.JsonPartialData{}, err
		}
		return termites_web.JsonPartialData{Key: "routing_graph", Data: msg}, nil
	},
)
