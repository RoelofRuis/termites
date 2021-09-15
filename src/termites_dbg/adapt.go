package termites_dbg

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/termites"
	"github.com/RoelofRuis/termites/termites_web"
)

var VisualizerAdapter = termites.NewAdapter(
	"Visualizer Path",
	"",
	termites_web.JsonPartialData{},
	func(i interface{}) (interface{}, error) {
		routingPath := i.(string)
		msg, err := json.Marshal(struct {
			RoutingPath string `json:"routing_path"`
		}{
			RoutingPath: routingPath,
		})
		if err != nil {
			return nil, err
		}
		return termites_web.JsonPartialData{Key: "visualizer", Data: msg}, nil
	},
)
