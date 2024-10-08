package termites_dbg

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites_web"
	"path/filepath"
)

var VisualizerAdapter = func(visualizerPath string) (termites_web.StateMessage, error) {
	_, filename := filepath.Split(visualizerPath)

	msg, err := json.Marshal(struct {
		Path string `json:"path"`
	}{
		Path: filepath.Join("/dbg-static/", filename),
	})
	if err != nil {
		return termites_web.StateMessage{}, err
	}
	return termites_web.StateMessage{Key: "graph", Data: msg}, nil
}
