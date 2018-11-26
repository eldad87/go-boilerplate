package logrus

import (
	log "github.com/sirupsen/logrus"
)

// Logger is an adapter from Jaeger Logger to Logrus.
type Logger struct {
	logger *log.Logger
}

// NewLogger creates a new Logger.
func NewLogger(logger *log.Logger) *Logger {
	return &Logger{logger: logger}
}

// Error logs a message at error priority
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Infof logs a message at info priority
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}
