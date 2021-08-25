package web

import (
	"bytes"
	"encoding/json"
	"github.com/RoelofRuis/termites/termites-core"
	"log"
)

type JsonCombiner struct {
	JsonDataIn  *termites.InPort
	JsonDataOut *termites.OutPort

	combiner *combiner
}

func NewJsonCombiner() *JsonCombiner {
	builder := termites.NewBuilder("JSON Combiner")

	combiner := &JsonCombiner{
		JsonDataIn:  builder.InPort("Partial Data", JsonPartialData{}),
		JsonDataOut: builder.OutPort("Data", []byte{}),

		combiner: newCombiner(),
	}

	builder.OnRun(combiner.Run)

	return combiner
}

func (f *JsonCombiner) Run(_ termites.NodeControl) error {
	data, err := f.combiner.get()
	if err != nil {
		log.Printf("JSON error: %s", err.Error())
	} else {
		f.JsonDataOut.Send(data)
	}

	for {
		select {
		case m := <-f.JsonDataIn.Receive():
			partialData := m.Data.(JsonPartialData)
			if !f.combiner.update(partialData) {
				continue
			}

			data, err := f.combiner.get()
			if err != nil {
				log.Printf("JSON error: %s", err.Error())
				continue
			}

			f.JsonDataOut.Send(data)
		}
	}
}

type JsonPartialData struct {
	Key  string
	Data json.RawMessage
}

type JsonData struct {
	Version int                        `json:"version"`
	Fields  map[string]json.RawMessage `json:"fields"`
}

type combiner struct {
	data JsonData
}

func newCombiner() *combiner {
	return &combiner{
		data: JsonData{
			Version: 0,
			Fields:  make(map[string]json.RawMessage),
		},
	}
}

func (c *combiner) get() ([]byte, error) {
	return json.Marshal(c.data)
}

func (c *combiner) update(data JsonPartialData) bool {
	existing, has := c.data.Fields[data.Key]
	if has && bytes.Compare(existing, data.Data) == 0 {
		return false
	}

	c.data.Fields[data.Key] = data.Data
	c.data.Version += 1

	return true
}
