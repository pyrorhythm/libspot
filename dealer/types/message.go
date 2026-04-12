package types

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	envTypeMessage = "message"
	envTypeRequest = "request"
)

// Envelope is the raw wire-level frame received over the dealer websocket.
// Both "message" and "request" frames share this shape; the Type field
// discriminates. Use IsMessage / IsRequest + ToMessage / ToRequest to convert
// into the typed variants that strip the fields that don't apply.
type Envelope struct {
	Type         string            `json:"type"`
	Method       string            `json:"method"`
	Uri          string            `json:"uri"`
	Headers      map[string]string `json:"headers"`
	MessageIdent string            `json:"message_ident"`
	Key          string            `json:"key"`
	Payloads     []json.RawMessage `json:"payloads"`
	Payload      compressed        `json:"payload"`
}

type compressed struct {
	Compressed []byte `json:"compressed"`
}

func (e *Envelope) IsMessage() bool { return e.Type == envTypeMessage }
func (e *Envelope) IsRequest() bool { return e.Type == envTypeRequest }

// ToMessage projects the envelope into a Message, dropping fields that only
// apply to requests (compressed payload, method).
func (e *Envelope) ToMessage() *Message {
	return &Message{
		Uri:          e.Uri,
		Headers:      e.Headers,
		MessageIdent: e.MessageIdent,
		Key:          e.Key,
		Payloads:     e.Payloads,
	}
}

// ToRequest projects the envelope into a Request, decompressing and
// decoding the gzipped JSON RequestPayload. Returns an error if decompression
// or decoding fails.
func (e *Envelope) ToRequest() (*Request, error) {
	if len(e.Payload.Compressed) == 0 {
		return nil, errors.New("dealer: request has no compressed payload")
	}
	gz, err := gzip.NewReader(bytes.NewReader(e.Payload.Compressed))
	if err != nil {
		return nil, fmt.Errorf("dealer: request gzip: %w", err)
	}
	defer gz.Close()
	raw, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("dealer: request gunzip: %w", err)
	}
	var p RequestPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("dealer: request json: %w", err)
	}
	return &Request{
		Key:     e.Key,
		Ident:   e.MessageIdent,
		Uri:     e.Uri,
		Method:  e.Method,
		Headers: e.Headers,
		Payload: p,
	}, nil
}

// Message is a fire-and-forget dealer frame. Payloads is per the wire spec
// always present; by convention the first element is the one consumers want.
type Message struct {
	Uri          string
	Headers      map[string]string
	MessageIdent string
	Key          string
	Payloads     []json.RawMessage
}
