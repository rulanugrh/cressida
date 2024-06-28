package helper

import (
	"fmt"
	"log/slog"
	"os"
)

type ILog interface {
	Info(message string, args ...any)
	Debug(message string)
	Error(err error)
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

func(l *Logger) Info(message string, args ...any) {
	msg := fmt.Sprintf("%v, %s", args, message)
	l.log.Info(msg)
}

func(l *Logger) Debug(message string) {
	l.log.Debug(message)
}

func(l *Logger) Error(err error) {
	l.log.Error(err.Error())
}

func(l *Logger) Warn(message string) {
	l.log.Warn(message)
}
