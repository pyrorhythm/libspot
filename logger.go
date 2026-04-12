package libspot

import "log/slog"

var _ StructuredLogger = (*slog.Logger)(nil)

type StructuredLogger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

func DefaultLogger() StructuredLogger {
	return slog.Default()
}