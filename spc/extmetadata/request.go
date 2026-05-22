package extmetadata

import (
	pb "pyrorhythm.dev/libspot/gen/spotify/extendedmetadata"
)

const Path = "extended-metadata/v0/extended-metadata"

type Request struct {
	country   string
	catalogue string
	taskID    []byte
	entities  []*pb.EntityRequest
}

func New() *Request { return &Request{} }

func (r *Request) Country(v string) *Request   { r.country = v; return r }
func (r *Request) Catalogue(v string) *Request { r.catalogue = v; return r }
func (r *Request) TaskID(v []byte) *Request    { r.taskID = v; return r }

func (r *Request) Query(uri string, kinds ...ExtensionKind) *Request {
	queries := make([]*pb.ExtensionQuery, len(kinds))
	for i, k := range kinds {
		queries[i] = pb.ExtensionQuery_builder{ExtensionKind: k}.Build()
	}
	return r.Entity(pb.EntityRequest_builder{
		EntityUri: uri,
		Query:     queries,
	}.Build())
}

func (r *Request) Entity(e *pb.EntityRequest) *Request {
	r.entities = append(r.entities, e)
	return r
}

func (r *Request) Build() *pb.BatchedEntityRequest {
	var header *pb.BatchedEntityRequestHeader
	if r.country != "" || r.catalogue != "" || r.taskID != nil {
		header = pb.BatchedEntityRequestHeader_builder{
			Country:   r.country,
			Catalogue: r.catalogue,
			TaskId:    r.taskID,
		}.Build()
	}

	return pb.BatchedEntityRequest_builder{
		Header:        header,
		EntityRequest: r.entities,
	}.Build()
}
