package zap

import (
	"context"

	"go.uber.org/zap"
)

// LoggerCreator represents a service for creating new *zap.Logger instance.
type LoggerCreator interface {
	// CreateLogger returns a new *zap.Logger instance.
	CreateLogger(ctx context.Context, componentName, methodName string) *zap.Logger
}

var _ LoggerCreator = (*LoggerZapCreator)(nil)

// LoggerZapCreator represents a service for creating new *zap.Logger instance.
type LoggerZapCreator struct {
	baseLogger *zap.Logger
}

// NewLoggerZapCreator returns a new LoggerZapCreator instance.
func NewLoggerZapCreator(baseLogger *zap.Logger) *LoggerZapCreator {
	return &LoggerZapCreator{
		baseLogger: baseLogger,
	}
}

// CreateLogger returns a new *zap.Logger instance.
func (creator *LoggerZapCreator) CreateLogger(ctx context.Context, componentName, methodName string) *zap.Logger {
	return creator.baseLogger.With(
		zap.String("component", componentName),
		zap.String("method", methodName))
}
