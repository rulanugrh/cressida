package helper

import (
	"log/slog"
	"os"
)

type ILog interface {
	Info(message string)
	Debug(message string)
	Error(err string)
	Warn(message string)
}

type Logger struct {
	log *slog.Logger
}

func NewLogger() ILog {
	return &Logger{
		log: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func(l *Logger) Info(message string) {
	l.log.Info(message)
}

func(l *Logger) Debug(message string) {
	l.log.Debug(message)
}

func(l *Logger) Error(err string) {
	l.log.Error(err)
}

func(l *Logger) Warn(message string) {
	l.log.Warn(message)
}
