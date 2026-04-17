package pathfinder

import (
	"encoding/json"

	"github.com/pyrorhythm/fn/bjs"
	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
	pfq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	pfs "github.com/pyrorhythm/libspot/pathfinder/pfresponse"
)

func AsPayload[T pfq.Request](t T) pfq.Payload[T] {
	return pfq.Payload[T]{
		Variables:     t,
		OperationName: t.Op(),
		Extensions:    &pfd.Extensions{PersistedQuery: t.Op().Extension()},
	}
}

func Marshal[T pfq.Request](t T) ([]byte, error) {
	return json.Marshal(AsPayload(t))
}

func Unmarshal[T any](data []byte) (*pfs.Payload[T], error) {
	return bjs.Unmarshal[pfs.Payload[T]](data)
}
