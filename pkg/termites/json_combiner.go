package termites

import (
	"bytes"
	"encoding/json"
)

type JsonCombiner struct {
	JsonDataIn  *InPort
	JsonDataOut *OutPort

	combiner *combiner
}

func NewJsonCombiner() *JsonCombiner {
	builder := NewBuilder("JSON Combiner")

	combiner := &JsonCombiner{
		JsonDataIn:  NewInPort[JsonPartialData](builder, "Partial Data"),
		JsonDataOut: NewOutPort[[]byte](builder, "Data"),

		combiner: newCombiner(),
	}

	builder.OnRun(combiner.Run)

	return combiner
}

func (f *JsonCombiner) Run(c NodeControl) error {
	data, err := f.combiner.get()
	if err != nil {
		c.LogError("JSON error", err)
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
				c.LogError("JSON error", err)
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
