package pfrequest

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type Payload[T any] struct {
	Variables     T               `json:"variables"`
	OperationName Operation       `json:"operationName"`
	Extensions    *pfd.Extensions `json:"extensions,omitempty"`
}
