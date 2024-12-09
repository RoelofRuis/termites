package termites_dbg

import (
	"encoding/json"
	"path/filepath"
)

type debuggerState struct {
	Graph struct {
		Path string `json:"path"`
	} `json:"graph"`
	Debugger struct {
		GraphEnabled    bool `json:"graph_enabled"`
		MessagesEnabled bool `json:"messages_enabled"`
		LogsEnabled     bool `json:"logs_enabled"`
	} `json:"debugger"`
}

func (d *debuggerState) Read() (json.RawMessage, error) {
	return json.Marshal(d)
}

type visualizerMessage struct {
	path string
}

func (v visualizerMessage) Mutate(s *debuggerState) error {
	_, filename := filepath.Split(v.path)
	s.Graph.Path = filepath.Join("/dbg-static/", filename)
	return nil
}
