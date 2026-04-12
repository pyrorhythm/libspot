package transport

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// LoggingTransport is an http.RoundTripper that logs request and response bodies.
type LoggingTransport struct{}

func (s *LoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	id := uuid.Must(uuid.NewV7())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	slog.Debug("request",
		"id", id.String(),
		"dest", r.URL.String(),
		"body", string(body),
	)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	slog.Debug("response",
		"id", id.String(),
		"code", resp.StatusCode,
		"body", string(respBody),
	)

	return resp, nil
}
