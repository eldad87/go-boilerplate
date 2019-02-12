package logger

import "go.uber.org/zap"

// Logger is an adapter from zap Logger to sql Logger.
type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger creates a new Logger.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger.Sugar()}
}

// Print is like fmt.Print
func (l *Logger) Print(args ...interface{}) {
	l.logger.Info(args...)
}
