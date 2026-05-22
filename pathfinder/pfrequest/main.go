package pfrequest

import pfd "pyrorhythm.dev/libspot/pathfinder/pfdomain"

type Payload[T Request] struct {
	Variables     T               `json:"variables"`
	OperationName Operation       `json:"operationName"`
	Extensions    *pfd.Extensions `json:"extensions,omitempty"`
}

type Request interface {
	Op() Operation
}
