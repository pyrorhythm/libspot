package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type Payload[T any] struct {
	Data       *T              `json:"data"`
	Extensions *pfd.Extensions `json:"extensions,omitempty"`
}
