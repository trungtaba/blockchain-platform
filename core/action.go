package core

import "encoding/json"

// Action ...
type Action struct {
	ActionName string          `json:"name"`
	Data       json.RawMessage `json:"data"`
}
