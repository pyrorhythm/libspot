package libspot

type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type NullLogger struct{}

func (l *NullLogger) Trace(...interface{}) {}
func (l *NullLogger) Debug(...interface{}) {}
func (l *NullLogger) Info(...interface{})  {}
func (l *NullLogger) Warn(...interface{})  {}
func (l *NullLogger) Error(...interface{}) {}
