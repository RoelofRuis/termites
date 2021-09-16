package termites_dbg

import (
	"encoding/json"
	"path/filepath"

	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_web"
)

var VisualizerAdapter = termites.NewAdapter(
	"Visualizer Path",
	"",
	termites_web.JsonPartialData{},
	func(i interface{}) (interface{}, error) {
		visualizerPath := i.(string)
		_, filename := filepath.Split(visualizerPath)

		msg, err := json.Marshal(struct {
			Path string `json:"path"`
		}{
			Path: filepath.Join("/dbg-static/", filename),
		})
		if err != nil {
			return nil, err
		}
		return termites_web.JsonPartialData{Key: "routing_graph", Data: msg}, nil
	},
)
