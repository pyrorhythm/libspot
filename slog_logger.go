package libspot

import (
	"fmt"
	"log/slog"
)

type slogLoggerAdapter struct {
	logger *slog.Logger
}

func NewSlogLogger(log *slog.Logger) Logger {
	return &slogLoggerAdapter{logger: log}
}

func (l *slogLoggerAdapter) Tracef(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *slogLoggerAdapter) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *slogLoggerAdapter) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *slogLoggerAdapter) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *slogLoggerAdapter) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *slogLoggerAdapter) Trace(args ...any) {
	l.logger.Debug(formatArgs(args...))
}

func (l *slogLoggerAdapter) Debug(args ...any) {
	l.logger.Debug(formatArgs(args...))
}

func (l *slogLoggerAdapter) Info(args ...any) {
	l.logger.Info(formatArgs(args...))
}

func (l *slogLoggerAdapter) Warn(args ...any) {
	l.logger.Warn(formatArgs(args...))
}

func (l *slogLoggerAdapter) Error(args ...any) {
	l.logger.Error(formatArgs(args...))
}

func (l *slogLoggerAdapter) WithField(key string, value any) Logger {
	return &slogLoggerAdapter{logger: l.logger.With(key, value)}
}

func (l *slogLoggerAdapter) WithError(err error) Logger {
	return &slogLoggerAdapter{logger: l.logger.With("error", err)}
}

func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			return s
		}
	}
	return fmt.Sprint(args...)
}
