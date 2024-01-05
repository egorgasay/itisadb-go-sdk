package itisadb

import "encoding/json"

type Value struct {
	Value    string
	ReadOnly bool
	Level    Level
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value    string
		ReadOnly bool
		Level    string
	}{
		Value:    v.Value,
		ReadOnly: v.ReadOnly,
		Level:    v.Level.String(),
	})
}
