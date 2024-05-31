package application

import (
	"fmt"
	"log/slog"
	"runtime"
)

type Logger struct {
	loggerInternal *slog.Logger
}

func NewLogger(internalLogger *slog.Logger) *Logger {
	return &Logger{loggerInternal: internalLogger}
}

func generateMessageWithError(message, location string, err error) string {
	return fmt.Sprintf("<%s>: %s Error: [%s]", location, message, err)
}

func generateMessageNoError(message, location string) string {
	return fmt.Sprintf("<%s>: %s", location, message)
}

func getCallerInfo() string {
	pc, _, _, _ := runtime.Caller(3)
	return runtime.FuncForPC(pc).Name()
}

func generateMessage(message string, err error) string {
	location := getCallerInfo()
	if err == nil {
		return generateMessageNoError(message, location)
	}
	return generateMessageWithError(message, location, err)
}

func (l *Logger) Debug(message string, err error, args ...any) {
	l.loggerInternal.Debug(generateMessage(message, err), args...)
}

func (l *Logger) Info(message string, err error, args ...any) {
	l.loggerInternal.Info(generateMessage(message, err), args...)
}

func (l *Logger) Warn(message string, err error, args ...any) {
	l.loggerInternal.Warn(generateMessage(message, err), args...)
}

func (l *Logger) Error(message string, err error, args ...any) {
	l.loggerInternal.Error(generateMessage(message, err), args...)
}
