package task

import (
	"encoding/json"
)

// Spec is the specification for a task
type Spec struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Init     string   `json:"init"`
	InitArgs []string `json:"initArgs"`
}

// MarshalBinary marshals a Spec
func (s *Spec) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary unmarshals a Spec
func (s *Spec) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// Info : task information
type Info struct {
	ID       string      `json:"id"`
	Metadata interface{} `json:"metadata"`
}
