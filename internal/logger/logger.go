package logger

import (
	"go.uber.org/zap"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	zapLogger *zap.Logger
}

// NewLogger initializes a new zap logger
func NewLogger() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{zapLogger: logger}
}

// Info logs informational messages
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.zapLogger.Info(msg, zapFields...)
}

// Error logs errors
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.zapLogger.Error(msg, zapFields...)
}

// Convert map fields to zap fields
func convertFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
