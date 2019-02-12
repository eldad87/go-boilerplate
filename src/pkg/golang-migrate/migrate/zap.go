package migrate

import "go.uber.org/zap"

// Logger is an adapter from zap Logger to migrate Logger.
type Logger struct {
	logger  *zap.SugaredLogger
	verbose bool
}

// NewLogger creates a new Logger.
func NewLogger(logger *zap.Logger, verbose bool) *Logger {
	return &Logger{logger: logger.Sugar(), verbose: verbose}
}

// Printf is like fmt.Printf
func (l *Logger) Printf(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Verbose should return true when verbose logging output is wanted
func (l *Logger) Verbose() bool {
	return l.verbose
}
