package pfrequest

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type Payload[T Request] struct {
	Variables     T               `json:"variables"`
	OperationName Operation       `json:"operationName"`
	Extensions    *pfd.Extensions `json:"extensions,omitempty"`
}

type Request interface {
	Op() Operation
}
