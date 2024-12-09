package examples

import "encoding/json"

// WebSharableState is an example struct implementing termites_web.State
type WebSharableState struct {
	Generator struct {
		Count int `json:"count"`
	} `json:"generator"`
}

func (w *WebSharableState) Read() (json.RawMessage, error) {
	return json.Marshal(w)
}
