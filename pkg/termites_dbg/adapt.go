package termites_dbg

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites_state"
	"path/filepath"
)

var VisualizerAdapter = func(visualizerPath string) (termites_state.StateMessage, error) {
	_, filename := filepath.Split(visualizerPath)

	msg, err := json.Marshal(struct {
		Path string `json:"path"`
	}{
		Path: filepath.Join("/dbg-static/", filename),
	})
	if err != nil {
		return termites_state.StateMessage{}, err
	}
	return termites_state.StateMessage{Key: "graph", Data: msg}, nil
}
