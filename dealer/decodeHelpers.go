package dealer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/pyrorhythm/libspot/dealer/types"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNoPayload           = errors.New("dealer: message has no payload")
	ErrTooManyPayloads     = errors.New("dealer: too many payloads in message")
	ErrUnknownDiscriminant = errors.New("dealer: unknown discriminant key")
)

// Predeclared topics. Each entry is the authoritative pairing of a dealer
// URI with the concrete Go type it carries. Adding a new built-in topic
// means adding one var here. User packages add their own by declaring
// their own Topic[T] values — no registry, no init-time wiring.
const connectionIDURIPrefix = "hm://pusher/v1/connections/"

// DecodePB is the standard TypedDecoder for payloads carrying a base64
// (optionally gzip-wrapped) protobuf in payloads[0]. Exported so user
// packages can compose their own Topic[T] without reimplementing the wire
// unwrap logic.
func DecodePB[T proto.Message](m *types.Message) (T, error) {
	var zero T
	if len(m.Payloads) == 0 {
		return zero, ErrNoPayload
	}
	if len(m.Payloads) > 1 {
		return zero, ErrTooManyPayloads
	}
	b, err := DecodeBytes(m.Payloads[0], m.Headers)
	if err != nil {
		return zero, err
	}
	msg := zero.ProtoReflect().New().Interface().(T)
	if err := proto.Unmarshal(b, msg); err != nil {
		return zero, fmt.Errorf("proto unmarshal: %w", err)
	}
	return msg, nil
}

// DecodeJSON is the standard TypedDecoder for payloads carrying an inline
// JSON object in payloads[0]. Returns a pointer to the decoded value.
func DecodeJSON[T any](m *types.Message) (*T, error) {
	if len(m.Payloads) == 0 {
		return nil, ErrNoPayload
	}
	if len(m.Payloads) > 1 {
		return nil, ErrTooManyPayloads
	}
	var v T
	if err := json.Unmarshal(m.Payloads[0], &v); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return &v, nil
}

// DecodeBytes unwraps a raw payload element into protobuf-ready bytes.
// The element is expected to be a JSON string containing base64 data,
// optionally gzipped when headers advertise Transfer-Encoding or
// Content-Encoding of "gzip". Exported for custom topics that decode into
// non-protobuf binary formats.
func DecodeBytes(raw json.RawMessage, headers map[string]string) ([]byte, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("payload not a string: %w", err)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64: %w", err)
	}
	if isGzip(headers) {
		gz, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("gzip reader: %w", err)
		}
		defer gz.Close()
		b, err = io.ReadAll(gz)
		if err != nil {
			return nil, fmt.Errorf("gzip unpack: %w", err)
		}
	}
	return b, nil
}

func isGzip(h map[string]string) bool {
	if h == nil {
		return false
	}
	for k, v := range h {
		if strings.EqualFold(k, "Transfer-Encoding") &&
			strings.Contains(strings.ToLower(v), "gzip") {
			return true
		}
		if strings.EqualFold(k, "Content-Encoding") &&
			strings.Contains(strings.ToLower(v), "gzip") {
			return true
		}
	}
	return false
}

func decodeConnectionID(m *types.Message) (string, error) {
	id := strings.TrimPrefix(m.Uri, connectionIDURIPrefix)
	if id == "" || id == m.Uri {
		return "", errors.New("dealer: empty connection id")
	}
	return id, nil
}

func decodeDeviceBroadcastStatus(m *types.Message) (*types.DeviceBroadcastStatus, error) {
	if len(m.Payloads) == 0 {
		return nil, ErrNoPayload
	}
	node, err := sonic.Get(m.Payloads[0])
	if err != nil {
		return nil, fmt.Errorf("sonic get: %w", err)
	}
	keymap, err := node.Map()
	if err != nil {
		return nil, fmt.Errorf("sonic map: %w", err)
	}
	if len(keymap) != 1 {
		return nil, fmt.Errorf("dealer: expected 1 discriminant key, got %d", len(keymap))
	}
	for k := range keymap {
		if k != "deviceBroadcastStatus" {
			return nil, fmt.Errorf("%w: %s", ErrUnknownDiscriminant, k)
		}
	}
	sub := node.Get("deviceBroadcastStatus")
	buf, err := sub.Raw()
	if err != nil {
		return nil, err
	}
	var v types.DeviceBroadcastStatus
	if err := json.Unmarshal([]byte(buf), &v); err != nil {
		return nil, err
	}
	return &v, nil
}
